CREATE TABLE `users` (
                         `id` varchar(255) NOT NULL,
                         `name` varchar(255) NOT NULL,
                         `email` varchar(255) NOT NULL,
                         `password` varchar(255) NOT NULL,
                         `role` enum(''Householder'',''ServiceProvider'',''Admin'') DEFAULT NULL,
                         `address` text,
                         `contact` varchar(20) DEFAULT NULL,
                         `is_active` tinyint(1) DEFAULT NULL,
                         `security_answer` varchar(255) DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `services` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` double DEFAULT NULL,
  `provider_id` varchar(255) DEFAULT NULL,
  `category` varchar(255) DEFAULT NULL,
  `avg_rating` double DEFAULT ''0'',
  `rating_count` int DEFAULT ''0'',
  UNIQUE KEY `unique_service_provider_category` (`provider_id`,`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `service_requests` (
                                    `id` varchar(255) NOT NULL,
                                    `householder_id` varchar(255) DEFAULT NULL,
                                    `householder_name` varchar(255) DEFAULT NULL,
                                    `householder_address` text,
                                    `service_id` varchar(255) DEFAULT NULL,
                                    `requested_time` datetime NOT NULL,
                                    `scheduled_time` datetime NOT NULL,
                                    `status` enum('Pending','Approved','Cancelled','Accepted') DEFAULT NULL,
                                    `approve_status` tinyint(1) DEFAULT '0',
                                    `service_name` varchar(255) DEFAULT NULL,
                                    `description` varchar(255) DEFAULT NULL,
                                    `householder_contact` varchar(255) DEFAULT NULL,
                                    PRIMARY KEY (`id`),
                                    KEY `householder_id` (`householder_id`),
                                    KEY `service_requests_ibfk_2` (`service_id`),
                                    CONSTRAINT `service_requests_ibfk_1` FOREIGN KEY (`householder_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `service_providers_services` (
                                              `service_provider_id` varchar(255) NOT NULL,
                                              `service_id` varchar(255) NOT NULL,
                                              `avg_rating` double DEFAULT '0',
                                              `rating_count` int DEFAULT '0',
                                              PRIMARY KEY (`service_provider_id`,`service_id`),
                                              KEY `service_id` (`service_id`),
                                              CONSTRAINT `service_providers_services_ibfk_1` FOREIGN KEY (`service_provider_id`) REFERENCES `service_providers` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `service_providers` (
                                     `user_id` varchar(255) NOT NULL,
                                     `rating` double DEFAULT '0',
                                     `availability` tinyint(1) DEFAULT '1',
                                     `is_active` tinyint(1) DEFAULT '1',
                                     PRIMARY KEY (`user_id`),
                                     CONSTRAINT `service_providers_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `service_provider_details` (
                                            `id` int NOT NULL AUTO_INCREMENT,
                                            `service_request_id` varchar(255) DEFAULT NULL,
                                            `service_provider_id` varchar(255) DEFAULT NULL,
                                            `name` varchar(255) DEFAULT NULL,
                                            `contact` varchar(20) DEFAULT NULL,
                                            `address` text,
                                            `price` varchar(50) DEFAULT NULL,
                                            `rating` double DEFAULT NULL,
                                            `approve` tinyint(1) DEFAULT NULL,
                                            `service_id` varchar(255) DEFAULT NULL,
                                            PRIMARY KEY (`id`),
                                            KEY `service_request_id` (`service_request_id`),
                                            KEY `service_provider_id` (`service_provider_id`),
                                            CONSTRAINT `service_provider_details_ibfk_1` FOREIGN KEY (`service_request_id`) REFERENCES `service_requests` (`id`),
                                            CONSTRAINT `service_provider_details_ibfk_2` FOREIGN KEY (`service_provider_id`) REFERENCES `service_providers` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9922 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `service_categories` (
                                      `id` varchar(255) NOT NULL,
                                      `category_name` varchar(255) NOT NULL,
                                      `description` varchar(255) DEFAULT NULL,
                                      PRIMARY KEY (`id`),
                                      UNIQUE KEY `category_name` (`category_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `reviews` (
                           `id` varchar(255) NOT NULL,
                           `service_id` varchar(255) DEFAULT NULL,
                           `householder_id` varchar(255) DEFAULT NULL,
                           `rating` double NOT NULL,
                           `comments` text,
                           `review_date` datetime DEFAULT CURRENT_TIMESTAMP,
                           `provider_id` varchar(255) DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           KEY `service_id` (`service_id`),
                           KEY `householder_id` (`householder_id`),
                           CONSTRAINT `reviews_ibfk_2` FOREIGN KEY (`householder_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
