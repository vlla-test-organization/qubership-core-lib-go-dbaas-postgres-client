package pgdbaas

import (
	"context"
	"strings"

	dbaasbase "github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/cache"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-postgres-client/v4/model"
)

const (
	DB_TYPE = "postgresql"
)

type Database interface {
	GetPgClient(options ...*model.PgOptions) (PgClient, error)
	GetConnectionProperties(ctx context.Context) (*model.PgConnProperties, error)
	FindConnectionProperties(ctx context.Context) (*model.PgConnProperties, error)
}

type database struct {
	params        model.DbParams
	dbaasPool     *dbaasbase.DbaaSPool
	postgresCache *cache.DbaaSCache
}

func (d database) GetPgClient(options ...*model.PgOptions) (PgClient, error) {
	clientOptions := &model.PgOptions{}
	if options != nil {
		clientOptions = options[0]
	}
	return &pgClientImpl{
		sqlOptions:      clientOptions.SqlOptions,
		bunOptions:      clientOptions.BunOptions,
		dbaasClient:     d.dbaasPool.Client,
		postgresqlCache: d.postgresCache,
		params:          d.params,
	}, nil
}

func (d database) GetConnectionProperties(ctx context.Context) (*model.PgConnProperties, error) {
	baseDbParams := d.params.BaseDbParams
	classifier := d.params.Classifier(ctx)

	pgLogicalDb, err := d.dbaasPool.GetOrCreateDb(ctx, DB_TYPE, classifier, baseDbParams)
	if err != nil {
		logger.Error("Error acquiring connection properties from DBaaS: %v", err)
		return nil, err
	}
	pgConnectionProperties := toPgConnProperties(pgLogicalDb.ConnectionProperties)
	return &pgConnectionProperties, nil
}

func (d database) FindConnectionProperties(ctx context.Context) (*model.PgConnProperties, error) {
	classifier := d.params.Classifier(ctx)
	params := d.params.BaseDbParams
	responseBody, err := d.dbaasPool.GetConnection(ctx, DB_TYPE, classifier, params)
	if err != nil {
		logger.ErrorC(ctx, "Error finding connection properties from DBaaS: %v", err)
		return nil, err
	}
	logger.Info("Found connection to pg db with classifier %+v", classifier)
	pgConnProperties := toPgConnProperties(responseBody)
	return &pgConnProperties, err
}

func toPgConnProperties(connProperties map[string]interface{}) model.PgConnProperties {
	url := connProperties["url"].(string)
	if strings.HasPrefix(url, "jdbc:") {
		url = strings.ReplaceAll(url, "jdbc:", "")
	}
	if tls, ok := connProperties["tls"].(bool); ok && tls {
		if tlsNotStrict, ok := connProperties["tlsNotStrict"].(bool); ok && tlsNotStrict {
			url += "?sslmode=require"
		} else {
			url += "?sslmode=verify-full"
		}
	}
	pgProp := model.PgConnProperties{
		Url:      url,
		Username: connProperties["username"].(string),
		Password: connProperties["password"].(string),
		Role:     connProperties["role"].(string),
	}
	if connProperties["roHost"] != nil {
		pgProp.RoHost = connProperties["roHost"].(string)
	}
	return pgProp
}
