<p align="center">
  <img src="https://user-images.githubusercontent.com/25181517/192149581-88194d20-1a37-4be8-8801-5dc0017ffbbe.png" width="100">
</p>
<h1 align="center">AvitoBTA2023-user_segmentation_service</h1>
<h3 align="center">2023 Entrance Tests for future <a href="https://avito.tech/">«Avito.Tech»</a> interns at <a href="https://www.avito.ru/company">«Avito»</a></h3>
<p align="center">The service provides HTTP API for managing user segments, including their creation, deletion and dynamic assignment to users taking into account the history of changes and TTL. It supports storing up-to-date data, fast updating of information and the ability to export change history to CSV.</p>

> [!IMPORTANT]
> The service was written hastily within 2-3 days. So, some things in the code are not finalized, but the whole functionality works correctly.

---

### — _Technology:_
![PostgreSQL](https://img.shields.io/badge/postgreSQL-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)

---

### — _How to Install and Use:_
#### 🟢 **Git Clone:**
```
git clone https://github.com/DigiRon4ik/AvitoBTA2023-user_segmentation_service.git
cd AvitoBTA2023-user_segmentation_service
```

> [!TIP]
> If you don't change the source code, when the service is raised, the builder container will work with cache, skipping capacitive processes.

#### 🟢 **Docker | Start/Stop:**
- 🚀
  ```
  docker-compose up -d
  ```
- ⛔
  ```
  docker-compose down
  ```
#### 🟢 **Make | Start/Stop:**
- 🚀
  ```
  make up
  ```
- ⛔
  ```
  make down
  ```

> [!NOTE]
> There is no migration mechanism, the database and tables are created by initializing the SQL script when the container is brought up.

---

### — _API Specification:_
1. The developed API is implemented according to the REST API design.
2. The body of the request and response are passed in JSON format.
3. Implemented Swagger and Swagger UI support for easy API handling.

---

### — _APIs:_
- **[Postman](https://www.postman.com/downloads/)**
  - You can import a [_JSON-file_](avitobta2023-user_segmentation_service.postman_collection.json) from the _API_ into _Postman_ for testing.
- **[Swagger UI](https://swagger.io/solutions/api-design/)**
  - Or use the _Swagger UI_ at [localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) after starting the service.
- **[Miro](https://miro.com/)**
  - Or for Russian-speakers you can see clearly in the _Miro_ (_clickable img_)
<p align="center">
  <a href="https://miro.com/app/live-embed/uXjVLmR3zL4=/?moveToViewport=-733,1765,7428,9330&embedId=989809006064">
    <img width=400 src="https://i.imgur.com/nfXc81A.png" >
  </a>
</p>

<div align="center">

#### Segments:
| Name             |     Method | API                |                                 Body                                 |
|:-----------------|-----------:|:-------------------|:--------------------------------------------------------------------:|
| Get all segments |    **GET** | `/segments`        |                                  -                                   |
| Get segment      |    **GET** | `/segments/{slug}` |                                  -                                   |
| Add segment      |   **POST** | `/segments`        | `{"slug": "AVITO_OFFER", "description": "Awaited offer (Optional)"}` |
| Update segment   |    **PUT** | `/segments/{slug}` |            `{"description": "Accepted offer (Optional)"}`            |
| Delete segment   | **DELETE** | `/segments/{slug}` |                                  -                                   |

#### Users:
| Name          |     Method | API                |         Body          |
|:--------------|-----------:|:-------------------|:---------------------:|
| Get all users |    **GET** | `/users`           |          -            |
| Get user      |    **GET** | `/users/{id}`      |          -            |
| Add user      |   **POST** | `/users`           | `{"name": "Abdulla"}` |
| Update user   |    **PUT** | `/users/{id}`      | `{"name": "Hayato"}`  |
| Delete user   | **DELETE** | `/users/{id}`      |          -            |

#### User Segments:
| Name                     |     Method | API                    |                                                                                  Body                                                                                 |
|:-------------------------|-----------:|:-----------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
| Get active user segments |    **GET** | `/users/{id}/segments` |                                                                                   -                                                                                   |
| Update user segments     |  **PATCH** | `/users/{id}/segments` | `{ "add": [ {"slug": "AVITO_VOICE_MESSAGES", "expiration_time": "2025-02-02T15:04:05Z" }, { "slug": "AVITO_DISCOUNT_30" } ], "remove": [ "AVITO_PERFORMANCE_VAS" ] }` |

#### User Segments History:
| Name                 |  Method | API                                                   |                                    Body                                   |
|:---------------------|--------:|:------------------------------------------------------|:-------------------------------------------------------------------------:|
| Update user segments | **GET** | `/users/{id}/segments/history?year={int}&month={int}` | `{"url": "http://localhost:8080/reports/report_{id}_{year}_{month}.csv"}` |

</div>

<p align="center">
  <a href="https://i.imgur.com/uP3zxzv.png">
    <img width=800 src="https://i.imgur.com/uP3zxzv.png" >
  </a>
</p>

---

### — _DataBase Schema:_
<p align="center">
  <a href="https://i.imgur.com/975GmDM.png">
    <img width=800 src="https://i.imgur.com/975GmDM.png" >
  </a>
</p>

---

### — _Progress in the fulfillment of assigned tasks:_

<div align="center">

| Task    | Progress | Comment     |
|:--------|:--------:|:------------|
| **Segment creation method** | ✅ | - |
| **Segment deletion method** | ✅ | - |
| **Method for adding a user to a segment** | ✅ | If you try to add or delete segments from a user again - they will be skipped or updated |
| **Method for obtaining active user segments** | ✅ | Since additional task #2 has been completed, the active segments are those with a `expiration_time` that has not yet occurred |
| **Code coverage by tests** | ❌ | I decided to skip it (don't hit me hard) |
| **Swagger** | ✅ | Described comments under swagger for handlers so that docs `swag init -g cmd/app/main.go -o api` can be generated |
| **Additional task No. 1 (*history*)** | ✅ | - |
| **Additional task No. 2 (*TTL*)** | ✅ | Support for deadline setting has been implemented - when the deadline expires, querying active user segments will not return a segment with an expired deadline (but it does not implement automatic table cleanup) |
| **Additional task No. 3 (*percentage*)** | ❌ | - |

</div>

---

### — _Video:_
<p align="center">
  <a href="https://youtu.be/27ToZvGJTVY">
    <img width=800 src="https://i.imgur.com/bqTOsir.png" >
  </a>
</p>

---

### — _Further development:_
1. Make the migration mechanism
2. Cover the code with tests
3. Maximize input data validation

---
