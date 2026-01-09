# FastAPI - Modern Backend Development

Time: 2 hours | Focus: Building production-ready APIs with FastAPI + Pydantic V2

## Why FastAPI in 2025?

- **Performance**: On par with Node.js and Go (thanks to Starlette + uvicorn)
- **Type safety**: Built on Python type hints + Pydantic V2
- **Auto documentation**: OpenAPI (Swagger) and ReDoc generated automatically
- **Async native**: Built for async/await from the ground up
- **Developer experience**: Autocomplete everywhere, fewer bugs

## Quick Start

```bash
# Install
uv add fastapi uvicorn[standard]

# Create app.py
```

```python
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():
    return {"message": "Hello World"}

# Run: uvicorn app:app --reload
# Visit: http://localhost:8000
# Docs: http://localhost:8000/docs
```

## Request and Response Models with Pydantic V2

### Basic Models

```python
from pydantic import BaseModel, EmailStr, Field, ConfigDict
from datetime import datetime

class UserCreate(BaseModel):
    """Request model for creating users"""
    email: EmailStr  # Validates email format
    password: str = Field(min_length=8, max_length=100)
    name: str = Field(min_length=1, max_length=100)
    age: int | None = Field(default=None, ge=0, le=150)

class UserResponse(BaseModel):
    """Response model - never return passwords!"""
    id: int
    email: str
    name: str
    created_at: datetime
    
    # Pydantic V2 configuration
    model_config = ConfigDict(from_attributes=True)  # For ORM objects

# Using in endpoints
@app.post("/users", response_model=UserResponse, status_code=201)
async def create_user(user: UserCreate):
    # user is validated automatically
    # FastAPI converts to dict with user.model_dump()
    db_user = await db.create_user(**user.model_dump())
    return db_user  # Auto-serialized to UserResponse
```

### Advanced Validation

```python
from pydantic import BaseModel, field_validator, model_validator
from typing import Annotated

class UpdateUser(BaseModel):
    email: EmailStr | None = None
    name: str | None = None
    password: str | None = None
    
    # Field-level validator
    @field_validator('password')
    @classmethod
    def password_strength(cls, v: str | None) -> str | None:
        if v is None:
            return v
        if len(v) < 8:
            raise ValueError('Password must be at least 8 characters')
        if not any(c.isupper() for c in v):
            raise ValueError('Password must contain uppercase letter')
        return v
    
    # Model-level validator (access multiple fields)
    @model_validator(mode='after')
    def check_at_least_one_field(self):
        if not any([self.email, self.name, self.password]):
            raise ValueError('At least one field must be provided')
        return self

# Using Annotated for reusable constraints
PositiveInt = Annotated[int, Field(gt=0)]
Username = Annotated[str, Field(min_length=3, max_length=50, pattern=r'^[a-zA-Z0-9_]+$')]

class Product(BaseModel):
    name: str
    price: PositiveInt
    quantity: PositiveInt
```

### Nested Models and Lists

```python
class Address(BaseModel):
    street: str
    city: str
    country: str
    postal_code: str

class OrderItem(BaseModel):
    product_id: int
    quantity: int
    price: float

class Order(BaseModel):
    user_id: int
    items: list[OrderItem]  # List of nested models
    shipping_address: Address  # Nested model
    notes: str | None = None

# FastAPI handles deep validation automatically
@app.post("/orders")
async def create_order(order: Order):
    # order.items[0].product_id is fully validated
    total = sum(item.quantity * item.price for item in order.items)
    return {"order_id": 123, "total": total}
```

## Path Parameters and Query Parameters

```python
from fastapi import Query, Path
from typing import Annotated

# Path parameters (required)
@app.get("/users/{user_id}")
async def get_user(
    user_id: Annotated[int, Path(ge=1)]  # Must be >= 1
):
    return {"user_id": user_id}

# Query parameters with validation
@app.get("/users")
async def list_users(
    page: Annotated[int, Query(ge=1)] = 1,
    limit: Annotated[int, Query(ge=1, le=100)] = 10,
    sort_by: Annotated[str, Query(max_length=50)] = "created_at",
    order: Annotated[str, Query(pattern="^(asc|desc)$")] = "desc",
    search: str | None = None
):
    """
    Query: /users?page=2&limit=20&sort_by=name&order=asc&search=alice
    """
    return {
        "page": page,
        "limit": limit,
        "sort_by": sort_by,
        "order": order,
        "search": search
    }

# Multiple path parameters
@app.get("/users/{user_id}/posts/{post_id}")
async def get_user_post(user_id: int, post_id: int):
    return {"user_id": user_id, "post_id": post_id}

# Enum for query parameters
from enum import Enum

class SortOrder(str, Enum):
    asc = "asc"
    desc = "desc"

@app.get("/products")
async def list_products(order: SortOrder = SortOrder.desc):
    return {"order": order.value}
```

## Dependency Injection (The Secret Sauce)

### Basic Dependencies

```python
from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession

# Database session dependency
async def get_db() -> AsyncSession:
    async with AsyncSessionLocal() as session:
        yield session

# Use in endpoints
@app.get("/users/{user_id}")
async def get_user(
    user_id: int,
    db: AsyncSession = Depends(get_db)  # Injected automatically
):
    user = await db.get(User, user_id)
    return user

# Reusable query parameters
class Pagination(BaseModel):
    page: int = Query(1, ge=1)
    limit: int = Query(10, ge=1, le=100)
    
    @property
    def offset(self) -> int:
        return (self.page - 1) * self.limit

@app.get("/users")
async def list_users(
    pagination: Pagination = Depends(),
    db: AsyncSession = Depends(get_db)
):
    users = await db.execute(
        select(User).offset(pagination.offset).limit(pagination.limit)
    )
    return users.scalars().all()
```

### Authentication with Dependencies

```python
from fastapi import HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

security = HTTPBearer()

async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: AsyncSession = Depends(get_db)
) -> User:
    token = credentials.credentials
    
    # Verify JWT token
    payload = verify_jwt_token(token)
    if not payload:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Get user from database
    user = await db.get(User, payload["user_id"])
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    return user

# Protected endpoint
@app.get("/me")
async def get_current_user_profile(
    current_user: User = Depends(get_current_user)
):
    return current_user

# Dependency with dependencies (chaining)
async def get_admin_user(
    current_user: User = Depends(get_current_user)
) -> User:
    if current_user.role != "admin":
        raise HTTPException(status_code=403, detail="Admin access required")
    return current_user

@app.delete("/users/{user_id}")
async def delete_user(
    user_id: int,
    admin: User = Depends(get_admin_user),  # Requires admin
    db: AsyncSession = Depends(get_db)
):
    await db.delete(await db.get(User, user_id))
    await db.commit()
    return {"status": "deleted"}
```

## Error Handling

```python
from fastapi import HTTPException, status
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError

# Raise HTTP exceptions
@app.get("/users/{user_id}")
async def get_user(user_id: int, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"User {user_id} not found",
            headers={"X-Error": "UserNotFound"}
        )
    return user

# Custom exception handlers
class DatabaseError(Exception):
    def __init__(self, message: str):
        self.message = message

@app.exception_handler(DatabaseError)
async def database_exception_handler(request, exc: DatabaseError):
    return JSONResponse(
        status_code=500,
        content={"error": "Database error", "detail": exc.message}
    )

# Override default validation error handler
@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request, exc: RequestValidationError):
    return JSONResponse(
        status_code=422,
        content={
            "error": "Validation failed",
            "details": exc.errors()
        }
    )

# Global exception handler
@app.exception_handler(Exception)
async def global_exception_handler(request, exc: Exception):
    # Log the error
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={"error": "Internal server error"}
    )
```

## Middleware

```python
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from starlette.middleware.base import BaseHTTPMiddleware
import time

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["https://example.com"],  # Or ["*"] for development
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Gzip compression
app.add_middleware(GZipMiddleware, minimum_size=1000)

# Custom middleware (timing)
class TimingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        start_time = time.time()
        response = await call_next(request)
        process_time = time.time() - start_time
        response.headers["X-Process-Time"] = str(process_time)
        return response

app.add_middleware(TimingMiddleware)

# Request logging middleware
class LoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        logger.info(f"{request.method} {request.url.path}")
        try:
            response = await call_next(request)
            logger.info(f"Status: {response.status_code}")
            return response
        except Exception as e:
            logger.error(f"Request failed: {e}")
            raise
```

## Background Tasks

```python
from fastapi import BackgroundTasks

def send_email(email: str, message: str):
    """Blocking email sending (runs in threadpool)"""
    time.sleep(2)  # Simulate slow email service
    print(f"Email sent to {email}: {message}")

async def log_activity(user_id: int, action: str, db: AsyncSession):
    """Async background task"""
    await db.execute(
        insert(ActivityLog).values(user_id=user_id, action=action)
    )
    await db.commit()

@app.post("/users")
async def create_user(
    user: UserCreate,
    background_tasks: BackgroundTasks,
    db: AsyncSession = Depends(get_db)
):
    # Create user
    db_user = await create_user_in_db(user, db)
    
    # Schedule background tasks (run after response sent)
    background_tasks.add_task(send_email, user.email, "Welcome!")
    background_tasks.add_task(log_activity, db_user.id, "user_created", db)
    
    return db_user
```

## Request and Response Examples

```python
from typing import Any

class UserResponse(BaseModel):
    id: int
    email: str
    name: str
    
    model_config = ConfigDict(
        json_schema_extra={
            "examples": [
                {
                    "id": 1,
                    "email": "alice@example.com",
                    "name": "Alice Smith"
                }
            ]
        }
    )

@app.post(
    "/users",
    response_model=UserResponse,
    status_code=201,
    summary="Create a new user",
    description="Creates a new user account with the provided information",
    response_description="The created user object",
    tags=["users"]
)
async def create_user(user: UserCreate):
    """
    Create a user with:
    
    - **email**: valid email address
    - **password**: minimum 8 characters
    - **name**: user's full name
    """
    pass
```

## File Uploads

```python
from fastapi import UploadFile, File

@app.post("/upload")
async def upload_file(file: UploadFile = File(...)):
    """
    Upload a file
    """
    contents = await file.read()
    
    return {
        "filename": file.filename,
        "content_type": file.content_type,
        "size": len(contents)
    }

# Multiple files
@app.post("/upload-multiple")
async def upload_multiple(files: list[UploadFile] = File(...)):
    return [
        {"filename": f.filename, "size": len(await f.read())}
        for f in files
    ]

# With form data
from fastapi import Form

@app.post("/upload-with-metadata")
async def upload_with_metadata(
    file: UploadFile = File(...),
    description: str = Form(...),
    tags: list[str] = Form([])
):
    return {
        "filename": file.filename,
        "description": description,
        "tags": tags
    }
```

## WebSockets (Real-time Communication)

```python
from fastapi import WebSocket, WebSocketDisconnect

class ConnectionManager:
    def __init__(self):
        self.active_connections: list[WebSocket] = []
    
    async def connect(self, websocket: WebSocket):
        await websocket.accept()
        self.active_connections.append(websocket)
    
    def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)
    
    async def broadcast(self, message: str):
        for connection in self.active_connections:
            await connection.send_text(message)

manager = ConnectionManager()

@app.websocket("/ws/{user_id}")
async def websocket_endpoint(websocket: WebSocket, user_id: int):
    await manager.connect(websocket)
    try:
        while True:
            data = await websocket.receive_text()
            await manager.broadcast(f"User {user_id}: {data}")
    except WebSocketDisconnect:
        manager.disconnect(websocket)
        await manager.broadcast(f"User {user_id} left")
```

## Application Lifespan Events

```python
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up...")
    await init_database()
    await init_redis()
    
    yield  # Application runs
    
    # Shutdown
    print("Shutting down...")
    await close_database()
    await close_redis()

app = FastAPI(lifespan=lifespan)
```

## Complete Example: User Management API

```python
from fastapi import FastAPI, Depends, HTTPException, status
from pydantic import BaseModel, EmailStr
from sqlalchemy.ext.asyncio import AsyncSession

app = FastAPI(title="User Management API", version="1.0.0")

# Models
class UserCreate(BaseModel):
    email: EmailStr
    name: str
    password: str

class UserResponse(BaseModel):
    id: int
    email: str
    name: str
    
    model_config = ConfigDict(from_attributes=True)

class UserUpdate(BaseModel):
    name: str | None = None
    email: EmailStr | None = None

# Dependencies
async def get_db() -> AsyncSession:
    async with AsyncSessionLocal() as session:
        yield session

# Endpoints
@app.post("/users", response_model=UserResponse, status_code=201)
async def create_user(
    user: UserCreate,
    db: AsyncSession = Depends(get_db)
):
    # Hash password
    hashed = hash_password(user.password)
    
    # Create user
    db_user = User(
        email=user.email,
        name=user.name,
        password_hash=hashed
    )
    db.add(db_user)
    await db.commit()
    await db.refresh(db_user)
    
    return db_user

@app.get("/users/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: int,
    db: AsyncSession = Depends(get_db)
):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user

@app.get("/users", response_model=list[UserResponse])
async def list_users(
    skip: int = 0,
    limit: int = 10,
    db: AsyncSession = Depends(get_db)
):
    result = await db.execute(
        select(User).offset(skip).limit(limit)
    )
    return result.scalars().all()

@app.patch("/users/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: int,
    user_update: UserUpdate,
    db: AsyncSession = Depends(get_db)
):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    # Update only provided fields
    update_data = user_update.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        setattr(user, field, value)
    
    await db.commit()
    await db.refresh(user)
    return user

@app.delete("/users/{user_id}", status_code=204)
async def delete_user(
    user_id: int,
    db: AsyncSession = Depends(get_db)
):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    await db.delete(user)
    await db.commit()
    return None
```

## Key Takeaways

1. **Pydantic models** for request/response validation
2. **Dependency injection** for database, auth, shared logic
3. **Type hints everywhere** - FastAPI uses them for validation
4. **Async/await native** - use async database drivers
5. **HTTPException** for error responses
6. **Auto-generated docs** at /docs and /redoc
7. **Background tasks** for post-response operations

Next: [04-database-sqlalchemy.md](./04-database-sqlalchemy.md) - Database operations with SQLAlchemy 2.0
