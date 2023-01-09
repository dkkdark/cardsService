# cardsService

Бэкенд для приложения: https://github.com/dkkdark/CardsServicesApp

Приложение что-то наподобии profi.ru (заказчики-исполнители)

Использована библиотека gorm. Для создания базы данных использовался PostgreSQL. 

Сама суть находится в файлах handlers и public. Каждый метод в handlers обращается к методам из publlic. Они выполняют select-запросы, update, insert, обращаются к функциям и процедурам БД.


Экраны регистрации и входа. Обращение к методам из public: AddUser, CheckUser

![Screenshot_20230109-112808_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211270424-d33c38b1-6efb-4921-b8e3-3078a50df4d2.jpg)
![Screenshot_20230109-112815_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211272560-61d3954e-9f1b-4e83-a8f1-9a51d8759430.jpg)


Экран с заданиями. Используется метод GetCards

![Screenshot_20230109-112940_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211273415-e7e646cd-94eb-4d99-aced-a3860ff57665.jpg)


Экран с информацией о заданиях и возможностью их бронировать. Используются медоты UpdateBookDatesUser, GetUserById

![Screenshot_20230109-113445_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211272851-8ab3d5ab-91fb-4bc2-a0a3-3f3d9f190a54.jpg)
![Screenshot_20230109-112954_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211272636-298c6b22-a9fc-40eb-90c9-ba789a6e5bd0.jpg)


Экран добавления здания. AddCard

![Screenshot_20230109-113020_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211273537-5cb5b2b6-2a00-47ea-a86d-863990b281a7.jpg)

Экран бронированных заданий. GetBookedCards

![Screenshot_20230109-113005_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211273587-6f357e0b-5b0e-4214-89c4-5102554a954e.jpg)


Экран заполнения профиля. GetSpecializationById, GetAddInfById, UpdateSpec, UpdateAddInf, UpdateCreatorStatus

![Screenshot_20230109-113215_Tasks App_photo-resizer ru](https://user-images.githubusercontent.com/49618961/211273733-b5bb6633-5dda-4af2-bcc0-da3235669f9a.jpg)


Есть и другие экраны, где используются остальные методы.
