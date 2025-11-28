package modify

func commands() string {
	return `package bootstrap

import "github.com/rusmanplatd/goravelframework/contracts/console"

func Commands() []console.Command {
	return []console.Command{}
}
`
}

func migrations() string {
	return `package bootstrap

import "github.com/rusmanplatd/goravelframework/contracts/database/schema"

func Migrations() []schema.Migration {
	return []schema.Migration{}
}
`
}

func providers() string {
	return `package bootstrap

import "github.com/rusmanplatd/goravelframework/contracts/foundation"

func Providers() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{}
}
`
}

func seeders() string {
	return `package bootstrap

import "github.com/rusmanplatd/goravelframework/contracts/database/seeder"

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{}
}
`
}
