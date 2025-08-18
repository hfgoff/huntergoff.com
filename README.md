See https://fly.io/docs/reference/configuration/ for information about how to use this file.

# articles that have informed the code in this repository

https://djwong.net/2025/05/28/cool-go-slog-tricks.html

# deps

go
brew install colima
brew install docker

fly deploy --build-arg COMMIT_HASH=$(git rev-parse HEAD) --build-arg VERSION=$(git describe --tags --always --dirty)

(rapid api)[https://rapidapi.com/movie-of-the-night-movie-of-the-night-default/api/streaming-availability/playground/apiendpoint_14b2f4b9-8801-499a-bcb7-698e550f9253]
