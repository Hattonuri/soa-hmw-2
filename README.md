## Сборка образов
Один раз нужно собрать образы:
```
    docker-compose build
```
## Запуск сервера
Боты имеют id в конце, чтобы их отличать.
```
    docker-compose up --scale bot=3 --scale player=0 --scale server=1
```

## Запуск реального игрока
```
    docker-compose run -e USER=Hattonuri player
```
Важно, присоединяться к уже запущенному серверу и
в другом окне терминала, чтобы можно было и играть за игрока и при желании переключаться
в окно с сервером и ботом для просмотра логов об игре.

## Работа с клиентом

Периодически происходит опрос(кого убить, вылечить) и достаточно ввести номер игрока для ответа
```
Time to vote...
Enter GetPlayers for get list players
bot531094794 (0)
bot1445028174 (1)
bot1504093019 (2)
Hattonuri (3)
```

Чтобы писать в общий чат - нужно как в minecraft писать ! перед текстом
Чтобы писать в чат мафии - нужно писать ?