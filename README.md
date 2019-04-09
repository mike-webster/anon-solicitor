# anon-solicitor

## Purpose
Provide a platform to faciliate soliciting public, anonymous feedback from a group of people.

## Requirements
- Allow users to create an account
    - Look into SSO (google)
- Allow authenticated users to see all submitted feedback
- Feedback will not have users attached to it
- Allow admin users the ability to configure settings
- Allow all authenticated users to create an "event" for which they would like public feedback
- Allow authenticated users to submit feedback for "events" to which they have been invited

## Dependencies
- Docker

## Endpoints
- **GET** `/events`
- **POST** `/events`
- **PUT** `/events/:id`
- **POST** `/events/:id/feedback`
- **GET** `/config`
- **PUT** `/config`
