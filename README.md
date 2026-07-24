# diplomas_crud

API CRUD para la primera fase de diplomas digitales.

Implementa persistencia base para:

- `documento_digital`
- `historico_estado_documento`

No implementa logica MID, integracion con firma electronica, SGA, S3 ni verificacion publica.

## Estado actual

Al viernes 24 de julio de 2026 el proyecto ya cuenta con:

- estructura base en Go y Beego v2;
- conexion ORM a PostgreSQL;
- migracion inicial del esquema `diplomas`;
- CRUD REST para `documento_digital` y `historico_estado_documento`;
- borrado logico con `activo = false`;
- filtros por `query`;
- base de CI/CD alineada en lo posible con `resoluciones_crud_v2`;
- contexto funcional y tecnico para iniciar el MID.

El flujo CRUD fue validado manualmente en Postman para altas, consultas, filtros, actualizaciones y borrado logico.

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

- [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
- [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
- [Docker](https://docs.docker.com/engine/install/ubuntu/)
- [Docker Compose](https://docs.docker.com/compose/)

## Variables de Entorno

```shell
DIPLOMAS_CRUD_HTTP_PORT=[Puerto de exposicion del API]
DIPLOMAS_CRUD_RUNMODE=[Modo de ejecucion: dev o prod]
DIPLOMAS_CRUD_PGUSER=[Usuario de BD]
DIPLOMAS_CRUD_PGPASS=[Contrasena del usuario de BD]
DIPLOMAS_CRUD_PGHOST=[URL, dominio o endpoint de la BD]
DIPLOMAS_CRUD_PGPORT=[Puerto de la BD]
DIPLOMAS_CRUD_PGDB=[Nombre de la base de datos]
DIPLOMAS_CRUD_PGSCHEMA=[Nombre del esquema]
DIPLOMAS_CRUD_PGSSLMODE=[Modo SSL de PostgreSQL]
DIPLOMAS_CRUD_PGTIMEZONE=[Zona horaria de la conexion]
PARAMETER_STORE=[Base path para secretos en AWS SSM, opcional]
```

**Nota:** Las variables se pueden ver en [conf/app.conf](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/conf/app.conf) y estan identificadas con `DIPLOMAS_CRUD_...`.

## Configuración local

Archivo `.env` base:

```env
DIPLOMAS_CRUD_HTTP_PORT=8080
DIPLOMAS_CRUD_RUNMODE=dev

DIPLOMAS_CRUD_PGHOST=127.0.0.1
DIPLOMAS_CRUD_PGPORT=5432
DIPLOMAS_CRUD_PGUSER=postgres
DIPLOMAS_CRUD_PGPASS=CHANGE_ME
DIPLOMAS_CRUD_PGDB=udistrital_academica
DIPLOMAS_CRUD_PGSCHEMA=diplomas
DIPLOMAS_CRUD_PGSSLMODE=disable
DIPLOMAS_CRUD_PGTIMEZONE=America/Bogota
```

Notas:

- `.env` se carga automaticamente al iniciar la app;
- si una variable ya existe exportada en shell, esa variable tiene prioridad;
- `conf/app.conf` funciona como fallback local.

## Ejecución del Proyecto

```shell
# 1. Entrar al repositorio
cd $GOPATH/src/github.com/udistrital/diplomas_crud

# 2. Ajustar variables en .env

# 3. Descargar dependencias
go mod tidy

# 4. Iniciar API
bee run
```

Alternativa:

```shell
go run .
```

## Ejecución Dockerfile

Implementado para despliegue CI/CD con binario `main` y `conf/app.conf`.

## Ejecución docker-compose

```shell
docker network create back_end
docker-compose up --build
```

## Migración inicial

Archivo:

[database/migrations/20260724_000001_init_diplomas.go](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/database/migrations/20260724_000001_init_diplomas.go)

La migracion crea:

- esquema `diplomas`;
- tabla `diplomas.documento_digital`;
- tabla `diplomas.historico_estado_documento`;
- indices y restricciones base de esta fase.

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

## Filtros soportados

`documento_digital`

```text
?query=hash_documento:abc
?query=uuid_verificacion:uuid
?query=tercero_id:123
?query=estado_documento_id:4
?query=activo:true
```

`historico_estado_documento`

```text
?query=documento_digital_id:10
?query=tercero_id:123
?query=estado_nuevo_id:8
?query=activo:true
```

Regla del parser:

- valores numericos se interpretan como enteros;
- `true` y `false` se interpretan como booleanos.

## Ejecución Pruebas

```shell
go test ./...
```

Si el entorno restringe cache Go:

```shell
GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...
```

## CI/CD

Archivos base:

- [Dockerfile](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/Dockerfile)
- [.drone.yml](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/.drone.yml)
- [sonar-project.properties](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/sonar-project.properties)
- [docker-compose.yml](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/docker-compose.yml)

Cobertura actual:

- `go vet`
- `go fmt`
- `golangci-lint`
- `go test ./...`
- build de binario `main`
- publicacion ECR para `release/*` y `master`
- actualizacion ECS para `release/*` y `master`
- Sonar preparado con secretos de Drone

Pendiente para despliegue real:

- registrar secretos `SONAR_HOST`, `SONAR_TOKEN`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`;
- confirmar nombre real del servicio ECS;
- decidir si se migra a stack institucional completo con `utils_oas`.

## Limitaciones actuales

- no hay autenticacion ni autorizacion;
- no hay Swagger generado;
- no hay pruebas HTTP automatizadas;
- el contrato final puede cambiar cuando DBA cierre tipos, restricciones o columnas;
- no existe aun logica de negocio para creacion automatica de historicos;
- no se integró `utils_oas` completo porque el modulo local actual depende de Beego v1 y este CRUD está en Beego v2.

## Archivos de contexto para el MID

- [CONTEXTO_API_MID.md](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/CONTEXTO_API_MID.md)
- [CONTEXTO_API_MID_DETALLADO.md](/home/vi/Documentos/SOPORTE/go/src/github.com/udistrital/diplomas_crud/CONTEXTO_API_MID_DETALLADO.md)
