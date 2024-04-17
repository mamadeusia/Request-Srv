CREATE TYPE "request_status" AS ENUM (
  'filling',
  'filled',
  'uncompleted',
  'admin_pending',
  'admin_checking',
  'admin_ask',
  'requester_admin_response',
  'validators_checking',
  'requester_validator_response',
  'blockchain_pending',
  'appoved',
  'rejected',
  'ValidatorQuestionReady'
);

CREATE TABLE "requests" (
  "id" UUID PRIMARY KEY,
  "full_name" varchar(128) NOT NULL,
  "age" int NOT NULL,
  "location_lat" float NOT NULL,
  "location_lon" float NOT NULL,
  "status" request_status NOT NULL,
  "requester_id" bigint UNIQUE NOT NULL,
  "photo" jsonb,
  "msgs" jsonb,
  "question_answers" jsonb,
  "created_at" timestamptz DEFAULT (now()) NOT NULL
);

CREATE TABLE "request_collaborators" (
  "request_id" UUID NOT NULL,
  "requester_id" bigint NOT NULL,
  "admin_id" bigint NOT NULL,
  "validators" jsonb
);

ALTER TABLE "request_collaborators" ADD FOREIGN KEY ("request_id") REFERENCES "requests" ("id");

ALTER TABLE "request_collaborators" ADD FOREIGN KEY ("requester_id") REFERENCES "requests" ("requester_id");



CREATE INDEX ON "requests" ("requester_id");
