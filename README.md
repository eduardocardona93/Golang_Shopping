# Golang Shopping API

API REST construida en Go con [Echo](https://echo.labstack.com/) y [GORM](https://gorm.io/) sobre PostgreSQL, siguiendo una arquitectura en capas (handler → service → repository) típica de proyectos Go en producción.

Gestiona cuatro entidades con CRUD completo:

- **Usuarios**
- **Productos**
- **Inventario**
- **Compras** (con control de stock disponible en Inventario)

## Arquitectura

```
cmd/api            → punto de entrada, wiring de dependencias, graceful shutdown
internal/config     → carga de configuración desde variables de entorno / .env
internal/database   → conexión GORM + PostgreSQL y auto-migración
internal/domain     → modelos GORM (entidades)
internal/dto        → request/response DTOs con validación (go-playground/validator)
internal/repository → acceso a datos (GORM), un repo por entidad
internal/service    → lógica de negocio (incluye la transacción de compra)
internal/handler    → controladores HTTP (Echo)
internal/router     → registro de rutas y middleware
pkg/apperrors       → errores de dominio (not found, conflict, validation, stock)
```

Cada capa depende solo de la interfaz de la capa inferior, lo que facilita testear con mocks y mantener la lógica de negocio libre de detalles HTTP o de base de datos.

## Modelo de datos

| Entidad | Campos |
|---|---|
| **User** (Usuario) | nombre, apellido, telefono, direccion, identificacion (única) |
| **Product** (Producto) | nombre, descripcion, sku (único), costo, precio, impuesto (% ) |
| **Inventory** (Inventario) | product_id (1:1 con Producto), cantidad, min_cantidad_alerta |
| **Purchase** (Compra) | user_id, items (lista de productos), fecha, impuesto (total), descuento_final, total |
| **PurchaseItem** (Lista de producto) | product_id, cantidad, precio_venta, descuento, impuesto |

Todas las entidades usan `id` tipo UUID y soft-delete (`deleted_at`).

### Reglas de negocio de las compras

- `precio_venta` e `impuesto` de cada línea se capturan del Producto en el momento de la compra (no se pueden falsear desde el cliente), garantizando trazabilidad histórica aunque el precio del producto cambie después.
- `descuento` (por línea) y `descuento_final` (por compra) son montos fijos, no porcentajes.
- Por cada línea: `subtotal = (precio_venta * cantidad - descuento) * (1 + impuesto / 100)`.
- `total = Σ subtotal - descuento_final`.
- **Control de stock**: al crear una compra, cada línea bloquea (`SELECT ... FOR UPDATE`) el registro de Inventario del producto dentro de una transacción, valida que `cantidad` solicitada no supere el stock disponible y lo decrementa atómicamente. Si el stock no alcanza, la API responde `409 Conflict` con el detalle del producto, cantidad solicitada y disponible, y no se aplica ningún cambio (rollback).
- **Actualizar una compra** libera el stock reservado por los ítems anteriores y vuelve a reservarlo con los nuevos ítems, dentro de la misma transacción.
- **Eliminar una compra** restaura el stock reservado antes de hacer soft-delete.

## Requisitos

- Go 1.25+
- Docker y Docker Compose (para levantar todo con un solo comando)
- PostgreSQL 16 (si se corre sin Docker)

## Levantar con Docker (recomendado)

```bash
cp .env.example .env
docker compose up --build
```

Esto levanta PostgreSQL y la API (puerto `8080` por defecto). La API espera a que la base de datos esté saludable y ejecuta la auto-migración al iniciar.

## Levantar en local (sin Docker)

```bash
cp .env.example .env   # ajustar credenciales si es necesario
go run ./cmd/api
```

## Variables de entorno

| Variable | Descripción | Default |
|---|---|---|
| `APP_ENV` | `development` habilita logs SQL detallados | `development` |
| `SERVER_PORT` | Puerto HTTP | `8080` |
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario | `postgres` |
| `DB_PASSWORD` | Password | `postgres` |
| `DB_NAME` | Base de datos | `shopping` |
| `DB_SSLMODE` | SSL mode de PostgreSQL | `disable` |

## Endpoints

Todos los endpoints de recursos están bajo el prefijo `/api/v1`. Formato de respuesta:

- Item: `{"data": {...}}`
- Lista: `{"data": [...], "meta": {"page", "limit", "total", "total_pages"}}`
- Error: `{"error": "mensaje", "fields": [...]?, "detail": {...}?}`

### Usuarios

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/users` | Crear usuario |
| GET | `/api/v1/users?page=&limit=` | Listar usuarios (paginado) |
| GET | `/api/v1/users/:id` | Obtener usuario |
| PUT | `/api/v1/users/:id` | Actualizar usuario |
| DELETE | `/api/v1/users/:id` | Eliminar usuario (soft delete) |

### Productos

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/products` | Crear producto |
| GET | `/api/v1/products?page=&limit=` | Listar productos |
| GET | `/api/v1/products/:id` | Obtener producto |
| PUT | `/api/v1/products/:id` | Actualizar producto |
| DELETE | `/api/v1/products/:id` | Eliminar producto |

### Inventario

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/inventories` | Crear registro de inventario para un producto |
| GET | `/api/v1/inventories?page=&limit=` | Listar inventario |
| GET | `/api/v1/inventories/:id` | Obtener inventario |
| PUT | `/api/v1/inventories/:id` | Actualizar cantidad / umbral de alerta |
| DELETE | `/api/v1/inventories/:id` | Eliminar registro de inventario |

### Compras

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/purchases` | Registrar compra (valida y descuenta stock) |
| GET | `/api/v1/purchases?page=&limit=` | Listar compras |
| GET | `/api/v1/purchases/:id` | Obtener compra con sus items |
| PUT | `/api/v1/purchases/:id` | Reemplazar items/fecha/descuento (reconcilia stock) |
| DELETE | `/api/v1/purchases/:id` | Eliminar compra (restaura stock) |

### Health check

`GET /health` → `{"status": "ok"}`

## Ejemplo: flujo completo de compra

```bash
# 1. Crear usuario
curl -X POST localhost:8080/api/v1/users -H 'Content-Type: application/json' -d '{
  "nombre": "Juan", "apellido": "Perez", "telefono": "3001234567",
  "direccion": "Calle 1", "identificacion": "CC123456"
}'

# 2. Crear producto
curl -X POST localhost:8080/api/v1/products -H 'Content-Type: application/json' -d '{
  "nombre": "Camiseta", "descripcion": "Camiseta de algodon", "sku": "SKU-001",
  "costo": 10.5, "precio": 25.0, "impuesto": 19
}'

# 3. Crear inventario para el producto
curl -X POST localhost:8080/api/v1/inventories -H 'Content-Type: application/json' -d '{
  "product_id": "<PRODUCT_ID>", "cantidad": 5, "min_cantidad_alerta": 2
}'

# 4. Registrar una compra (falla con 409 si la cantidad supera el stock disponible)
curl -X POST localhost:8080/api/v1/purchases -H 'Content-Type: application/json' -d '{
  "user_id": "<USER_ID>",
  "items": [{"product_id": "<PRODUCT_ID>", "cantidad": 3, "descuento": 5}],
  "descuento_final": 2
}'
```

## Comandos útiles (Makefile)

```bash
make run          # go run ./cmd/api
make build        # compila el binario en bin/api
make test         # go test ./...
make vet          # go vet ./...
make tidy         # go mod tidy
make docker-up    # docker compose up --build -d
make docker-down  # docker compose down
```
