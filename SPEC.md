## Product Definition
A seat plan application.
- Admins define seats in a room.
- Each room could have certatin number of seats.
- Each seat could be used by one person in a day.
- Users have a weekly reservation limit (e.g., 2 days per week). Team admins can configure this limit for their team members.
- We also have teams that owns certain seats.
- Each team has some member and members doesn't overlap.
- There is viewer users, users who could reserve seat for itself and admins who could edit seat plans for everyone.
- There is superadmin that has every permission and there is team admins that could everything in a team not for other teams.
- Team admins could request some seat for some day from other team admin, and if other admin accept other teams could reserve other team seats.
- We show each seat as a rectangle that has an axis.
- Each seat has a label too.
- Each room is a grid that could we change the size of it. 

## Tech Spec
- This project should be a monorepo containing at least two frontend application and one backend application.
- We are going to build the backend app using golang.
- We are going to develop the backend using hexagonal architecture.
- For database we are using postgres.
- We use sqlc on the backend service
- For the frontend we use vite, tailwindcss and typescript.
- There should be proper eslint, prettier config with it's git push constraints.
- We use solid js as the frontend framework.
- Every application in the monorepo should be dockerized.
- We use a docker compose for running everything together on the production.

