#!/bin/sh

su postgres
psql
CREATE USER master  WITH PASSWORD '12345';
CREATE DATABASE forum;

