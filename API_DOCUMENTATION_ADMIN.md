# Admin API Documentation

## 1. Create Team
- **URL**: `/api/teams/create`
- **Method**: `POST`
- **Headers**: 
  - `Authorization`: `Bearer <token>`
  - `Content-Type`: `application/json`
- **Request Body**:
  ```json
  {
    "name": "Team Name",
    "sport": "football",
    "college": "Computer Science",
    "team_type": "Men's"
  }
  ```
- **Response Example**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "id": 1,
      "name": "Team Name",
      ...
    }
  }
  ```

## 2. Create Event
- **URL**: `/api/events/create`
- **Method**: `POST`
- **Headers**: 
  - `Authorization`: `Bearer <token>`
  - `Content-Type`: `application/json`
- **Request Body**:
  ```json
  {
    "name": "Event Name",
    "sport": "football",
    "start_date": "2024-05-01",
    "format_type": "Knockout"
  }
  ```
- **Response Example**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "id": 1,
      "name": "Event Name",
      ...
    }
  }
  ```

## 3. Create Match
- **URL**: `/api/matches/create`
- **Method**: `POST`
- **Headers**: 
  - `Authorization`: `Bearer <token>`
  - `Content-Type`: `application/json`
- **Request Body**:
  ```json
  {
    "event_id": 1,
    "team_a_id": 1,
    "team_b_id": 2,
    "name": "Final Match",
    "time": "2024-05-02T14:00:00Z"
  }
  ```
- **Response Example**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "id": 1,
      "name": "Final Match",
      ...
    }
  }
  ```

## 4. List Teams
- **URL**: `/api/teams/list`
- **Method**: `GET`
- **Headers**: 
  - `Authorization`: `Bearer <token>`
- **Response Example**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": [
      {
        "id": 1,
        "name": "Team A",
        "sport": "football",
        "college": "CS",
        "team_type": "Men's"
      }
    ]
  }
  ```

## 5. List Events
- **URL**: `/api/events/list`
- **Method**: `GET`
- **Headers**: 
  - `Authorization`: `Bearer <token>`
- **Response Example**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": [
      {
        "id": 1,
        "name": "Event A",
        "sport": "football",
        "start_date": "2024-05-01",
        "format_type": "Knockout"
      }
    ]
  }
  ```
