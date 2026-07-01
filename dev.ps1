param(
    [ValidateSet("up", "rebuild", "logs", "ps", "down")]
    [string]$Action = "up",

    [ValidateSet("backend", "frontend", "nginx", "minio")]
    [string]$Service
)

$ErrorActionPreference = "Stop"

switch ($Action) {
    "up" {
        docker compose up -d --build
    }
    "rebuild" {
        if ($Service) {
            docker compose up -d --build --force-recreate --no-deps $Service
        } else {
            docker compose up -d --build --force-recreate
        }
    }
    "logs" {
        if ($Service) {
            docker compose logs -f $Service
        } else {
            docker compose logs -f
        }
    }
    "ps" {
        docker compose ps
    }
    "down" {
        docker compose down
    }
}
