# anon-solicitor

## Purpose
Provide a platform to faciliate soliciting public, anonymous feedback from a group of people.

## Requirements
- Display all feedback, organized by event.
    - Feedback will not have users attached to it
- Any user can anonymously set up an event.
    - Always ask the following:
        - 1-5 Hate - Love the event
        - What could we have done better?
        - Did anything go poorly?
    - Provide any other specific questions
    - Provide a list of emails for the audience
- Audience receives emails
    - The emails will contain a one-time use token allowing that user to provide feedback to the event.
    - The email will also contain a link that will inform the results that the user didn't attent the event.
        - This will also invalidate the one-time use token.

## Dependencies
- Docker

## Endpoints
- **GET** `/events`
- **POST** `/events`
- **PUT** `/events/:id`
- **POST** `/events/:id/feedback`
- **GET** `/config`
- **PUT** `/config`

## Flags
- `drop={bool}` will specify if you want to drop the tables when the app starts; defaults to false.