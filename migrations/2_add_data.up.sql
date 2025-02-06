INSERT INTO users (username) VALUES
('user1'),
('user2'),
('user3');

INSERT INTO posts (author_id, title, content, allow_comments) VALUES
(1, 'Первый пост', 'Это тестовый пост пользователя 1', true),
(2, 'Второй пост', 'Это тестовый пост пользователя 2', true);

INSERT INTO comments (post_id, author_id, content) VALUES
(1, 2, 'Комментарий пользователя 2 к посту 1'),
(1, 3, 'Комментарий пользователя 3 к посту 1'),
(2, 1, 'Комментарий пользователя 1 к посту 2');
