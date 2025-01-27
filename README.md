# Movie Reservation Service

## Overview

The Movie Reservation Service is a backend system designed to facilitate user authentication, movie management, seat reservation, and reporting for a movie reservation platform. This service allows users to sign up, log in, browse movies, reserve seats for specific movies. Admins have additional privileges to manage movies.

## Goals

The primary goal of this project is to implement complex business logic related to seat reservation and scheduling, while also understanding data models, relationships, and complex queries.

## Features

- **User Authentication and Authorization**
  - User sign-up and login functionality.
  - Role-based access control (admin and regular user).
  - Admins can promote users to admin and manage movies.
  - Use JWT for secure user authentication and role-based access control.

- **Movie Management**
  - Admins can add, update, and delete movies.
  - Movies have attributes such as title, description, poster image, and genre.
  - Movies are associated with showtimes.
  - Design a relational database schema to represent movies and their relationships.

- **Reservation Management**
  - Users can view movies and their showtimes for a specific date.
  - Users can reserve seats for a showtime and see available seats.
  - Users can view and cancel their upcoming reservations.
  - Admins can view all reservations, capacity, and revenue.
  - Implement logic to avoid overbooking and handle seat reservations effectively.
  - Manage the scheduling of showtimes to ensure accurate availability.

- **Reporting**
  - Provide reporting features for admins to view reservations, capacity, and revenue.


## Technologies Used

- **Programming Language:** Go
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Environment Management:** godotenv

