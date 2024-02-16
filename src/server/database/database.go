package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"thor/src/server/config"
)

type Database struct {
	*gorm.DB
}

func InitializeDatabase(ds config.Datasource) *Database {

	//mysql conn
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	ds.Username,
	//	ds.Password,
	//	ds.Url,
	//	ds.Port,
	//	ds.DatabaseName)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require search_path=%s",
		ds.Url,
		ds.Username,
		ds.Password,
		ds.DatabaseName,
		ds.Port,
		ds.Schema)

	cfg := &gorm.Config{}

	//rootCertPool := x509.NewCertPool()
	//pem, err := ioutil.ReadFile("/src/db_cert/ca.pem")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
	//	fmt.Println("Failed to append PEM.")
	//}
	//clientCert := make([]tls.Certificate, 0, 1)
	//
	//certs, err := tls.LoadX509KeyPair("src/db_cert/rootCACert.pem", "src/db_cert/ca.pem")
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//clientCert = append(clientCert, certs)
	//
	//cfg = &gorm.Config{tls.Certificate{Certificate: certs, PrivateKey: pem}}

	if ds.DebugMode {
		cfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	}

	cfg.NamingStrategy = schema.NamingStrategy{TablePrefix: ds.Schema, SingularTable: false}

	db, err := gorm.Open(postgres.Open(dsn), cfg)

	if err != nil {
		fmt.Println("Failed to open connection")
		panic(err)
	}

	sqlDb, err := db.DB()

	if err != nil {
		fmt.Println("No database found")
		panic(err)
	}

	sqlDb.SetConnMaxIdleTime(ds.ConnectionTimeout)
	sqlDb.SetMaxIdleConns(ds.MaxIdleConnection)
	sqlDb.SetMaxOpenConns(ds.MaxOpenConnection)

	return &Database{db}
}
