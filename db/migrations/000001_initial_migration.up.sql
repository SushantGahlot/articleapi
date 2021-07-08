CREATE TABLE IF NOT EXISTS "article" (
  "article_id" UUID NOT NULL UNIQUE PRIMARY KEY,
  "article_title" varchar(255) NOT NULL,
  "article_date" date NOT NULL,
  "article_body" text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "article_tag" (
  "article_id" UUID NOT NULL,
  "tag_id" UUID NOT NULL,
  PRIMARY KEY ("article_id", "tag_id")
);

CREATE TABLE IF NOT EXISTS "tags" (
  "tag_id" UUID NOT NULL UNIQUE PRIMARY KEY,
  "tag" varchar NOT NULL UNIQUE
);

ALTER TABLE "article_tag" ADD FOREIGN KEY ("article_id") REFERENCES "article" ("article_id");

ALTER TABLE "article_tag" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("tag_id");
