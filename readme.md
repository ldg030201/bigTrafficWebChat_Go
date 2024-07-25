# DDL

```SQL
CREATE DATABASE `bigtrafficwebchat` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;

CREATE TABLE `room` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '룸 고유 ID',
    `name` varchar(255) NOT NULL COMMENT '룸 이름',
    `createAt` timestamp NULL DEFAULT current_timestamp() COMMENT '만든 일시',
    `updateAt` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '수정된 일시',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='룸 데이터';

CREATE TABLE `chat` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '채팅 고유 ID',
    `room` varchar(255) NOT NULL COMMENT '룸 고유 ID',
    `name` varchar(255) NOT NULL COMMENT '송신자 이름',
    `message` varchar(255) NOT NULL COMMENT '메시지 내용',
    `when` timestamp NULL DEFAULT current_timestamp() COMMENT '메시지 보낸 시간',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='채팅 데이터';

CREATE TABLE `serverinfo` (
    `ip` varchar(255) NOT NULL COMMENT '서버 IP',
    `avaliable` tinyint(1) NOT NULL COMMENT '서버 온 오프 여부',
    PRIMARY KEY (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='서버정보';


ALTER TABLE room ADD CONSTRAINT room_unique UNIQUE KEY (name);
```