# diplomas_crud

API CRUD para la primera fase de diplomas digitales.

Implementa persistencia base para:

- `documento_digital`
- `historico_estado_documento`

## Estado CI

| Develop | Release 0.0.1 | Master | Sonar |
| --- | --- | --- | --- |
| Pendiente | Pendiente | Pendiente | Pendiente |

### Tecnologías Implementadas y Versiones

* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)


## Variables de entorno

```shell
DIPLOMAS_CRUD_HTTP_PORT=[Puerto de exposicion del API]
DIPLOMAS_CRUD_RUNMODE=[Modo de ejecucion: dev o prod]
DIPLOMAS_CRUD_PGUSER=[Usuario de BD]
DIPLOMAS_CRUD_PGPASS=[Contrasena del usuario de BD]
DIPLOMAS_CRUD_PGHOST=[Host de la BD]
DIPLOMAS_CRUD_PGPORT=[Puerto de la BD]
DIPLOMAS_CRUD_PGDB=[Nombre de la base de datos]
DIPLOMAS_CRUD_PGSCHEMA=[Nombre del esquema]
PARAMETER_STORE=[Base path para secretos en AWS SSM, opcional]
```

Configuración base en [conf/app.conf](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/conf/app.conf).

## Ejecución local

```shell
cd $GOPATH/src/github.com/udistrital/diplomas_crud
go get github.com/udistrital/utils_oas@latest
go mod tidy
bee run
```

Alternativa:

```shell
go run .
```

## Pruebas

```shell
go test ./...
```

Si el entorno restringe caché Go:

```shell
GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...
```

## Docker

```shell
docker network create back_end
docker-compose up --build
```

## Endpoints

```text
GET    /v1/documento_digital/
POST   /v1/documento_digital/
GET    /v1/documento_digital/:id
PUT    /v1/documento_digital/:id
DELETE /v1/documento_digital/:id

GET    /v1/historico_estado_documento/
POST   /v1/historico_estado_documento/
GET    /v1/historico_estado_documento/:id
PUT    /v1/historico_estado_documento/:id
DELETE /v1/historico_estado_documento/:id
```

## Licencia

This file is part of diplomas_crud.

diplomas_crud is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

diplomas_crud is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with diplomas_crud. If not, see [https://www.gnu.org/licenses/](https://www.gnu.org/licenses/).
