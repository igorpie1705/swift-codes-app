# Swift Codes API

## Wstęp
Ten projekt to API do zarządzania kodami SWIFT. Obsługuje operacje CRUD (Create, Read, Update, Delete) dla kodów SWIFT i przechowuje je w bazie danych PostgreSQL.

## Wymagania
- **Docker** i **Docker Compose** (do uruchomienia aplikacji)
- **Postman** (opcjonalnie, do testowania endpointów)
- **Git** (do pobrania repozytorium)
- **cURL** (opcjonalnie, do testowania API w terminalu)

## Instalacja i uruchomienie
### 1. Klonowanie repozytorium
Najpierw otwórz terminal (Linux/Mac) lub PowerShell (Windows) i wpisz:
```sh
git clone https://github.com/igorpie1705/swift-codes-app.git
cd swift-codes-app
```

### 2. Uruchomienie aplikacji w Dockerze
Jeśli masz już zainstalowanego Dockera, wpisz w terminalu:
```sh
docker-compose up --build
```
To uruchomi:
- PostgreSQL jako bazę danych
- API w kontenerze Go

Aplikacja będzie dostępna pod adresem **http://localhost:8080**.

## Struktura projektu
```
.
├── database/        # Obsługa połączenia z bazą danych
├── handlers/        # Logika biznesowa API
├── models/          # Definicje modeli
├── tests/           # Testy integracyjne i jednostkowe
├── main.go          # Punkt startowy aplikacji
└── docker-compose.yml  # Konfiguracja Dockera
```

## Endpointy API

### Pobranie kodu SWIFT
**GET /v1/swift-codes/:swift-code**
```sh
curl -X GET http://localhost:8080/v1/swift-codes/ABCDEFGXXX
```

### Pobranie kodów SWIFT dla kraju
**GET /v1/swift-codes/country/:countryISO2code**
```sh
curl -X GET http://localhost:8080/v1/swift-codes/country/PL
```

### Dodanie nowego kodu SWIFT
**POST /v1/swift-codes**
```sh
curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d '{
    "swiftCode": "NEWBANKGBXXX",
    "bankName": "New Bank",
    "address": "789 New St",
    "countryISO2": "GB",
    "countryName": "United Kingdom",
    "isHeadquarter": true
}'
```

### Usunięcie kodu SWIFT
**DELETE /v1/swift-codes/:swift-code**
```sh
curl -X DELETE http://localhost:8080/v1/swift-codes/NEWBANKGBXXX
```

## Testowanie API w Postman
1. Otwórz **Postman**
2. Dodaj nowe żądanie **GET/POST/DELETE**
3. Użyj URL **http://localhost:8080/v1/swift-codes**
4. Kliknij "Send" i sprawdź odpowiedź

## Uruchamianie testów
Testy można uruchomić ręcznie w terminalu:
```sh
go test ./tests -v
```

### Uruchamianie testów w Dockerze
Jeśli chcesz uruchomić testy w kontenerze Dockera, wpisz:
```sh
docker-compose run app go test ./tests -v
```
To sprawdzi poprawność działania API oraz połączenia z bazą danych.

## Zakończenie
Jeśli masz pytania lub napotkasz błędy, zapraszam do kontaktu! :)

