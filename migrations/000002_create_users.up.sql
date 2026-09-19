-- up
CREATE TABLE myapp.users (
                             id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                             name TEXT NOT NULL,
                             age INTEGER
);