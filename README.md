# Задание на стажировку VK - REST API сервис 

Инструкция по запуску:

1. С помощью флагов или переменных окружения при запуске сервера задать параметры конфигурации:
    1) порт - по умолчанию ":8080", флаг "a", переменная - "SERVER_ADDRESS".
    2) URL - по умолчанию "http://localhost:8080", флаг "b", переменная - "BASE_URL".
    3) путь и данные доступа к БД PostgreSQL (по умолчанию установлены локальные), флаг "d", переменная - "DATABASE_DSN".
2. Аналогичнм образом при запуске клиента задать URL (по умолчанию "http://localhost:8080", флаг "b", переменная - "BASE_URL").
3. Клиент выполнен в формате CLI-приложения, управление возможными действиями осуществляется путем ввода данных с клавиатуры. При возникновении ошибки (обычно при вводе неверных данных) необходимо повторно запустить клиент.

Сервис содержит следующие методы:
1. Добавление нового пользователя по имени - возвращается его ID.
2. Добавление нового задания по имени и стоимости - возвращается его ID.
3. Выполнение определенного задания конкретным пользователем (по ID) (возможно только 1 раз) - возвращается название выполеннного задания.
4. Просмотр баланса и истории выполнения заданий конкретным пользователем (по ID) - возвращается баланс и история пользователя.