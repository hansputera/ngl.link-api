
# Ngl.Clone API

I just want to make my own ngl.link, and learn how they works.

## Installation

You need to install Golang first! Then

```bash
  git clone https://github.com/hansputera/ngl.link-api
  cd ngl.link-api

  go mod download
```
    
## Run Locally

Build the source code

```bash
  go build
```

Run it

```bash
  ./nglapi
```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`JWT_SECRET` Fill it with random chars

`REDIS_URI` Fill it with your redis connection URI


## License

[MIT](https://choosealicense.com/licenses/mit/)

