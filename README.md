[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/slinkity/slinkity/blob/main/LICENSE.md)

# kubeground

> 🚧 **This project is heavily under construction!** 🚧

Discover kubeground, an open-source app that provides Kubernetes training in a unique and engaging way. With hand-picked challenges, you can learn Kubernetes by troubleshooting different components of the system. As you progress through the levels, you'll gain valuable experience and deepen your understanding of this powerful technology.

- 🎮 **Gamify kubernetes learning** Experience a fun and engaging way to improve your Kubernetes knowledge with our hand-picked, unique challenges. Solve each level by troubleshooting your local cluster, and enjoy the process of leveling up your skills.
- 🔋 **Batteries included** Have no lagging experience while using the product (thanks Echo!).
- 📖 **Comprehensive level documentation included** All levels should have clear description & documentation.

## Technologies used

kubeground uses many popular open source tools, including:

1. [**Go Echo:**](https://echo.labstack.com/) High performance, extensible, minimalist Go web framework.
2. [**GORM:**](https://gorm.io/) The fantastic ORM library for Golang, aims to be developer friendly.
3. [**kubernetes/client-go:**](https://github.com/kubernetes/client-go) Go clients for talking to a kubernetes cluster.
4. [**Kind:**](https://kind.sigs.k8s.io/) kind is a tool for running local Kubernetes clusters using Docker container “nodes”.
   kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.
5. [**Tilt:**](https://tilt.dev/) A toolkit for fixing the pains of microservice development.
6. [**Vite:**](https://vitejs.dev) a bundler that takes the boilerplate out of your set up. It'll compile JS component frameworks, handle CSS preprocessors with little-to-no config (say, SCSS and PostCSS), and show dev changes on-the-fly using [hot module replacement (HMR)](https://vitejs.dev/guide/features.html#hot-module-replacement).

7. [**React:**](https://reactjs.org/) React is a JavaScript library for building user interfaces.
8. [**MUI:**](https://mui.com/) MUI Core contains foundational React UI component libraries for shipping new features faster.

## Getting started

```go
cd backend
go run .
```

Navigate to the local [dashboard](localhost:4000/) to start init some levels

## Feature set

This project is still in early alpha, so we have many features soon to come!<br>
For reference, here's our complete roadmap of current and upcoming features:

| Feature                                                                                                         | Status |
| --------------------------------------------------------------------------------------------------------------- | ------ |
| Create and delete kubernetes resources dynamically from YAML files:                                             | ✅     |
| Build a fast API server with Go Echo to handle business logic                                                   | ✅     |
| Utilize GORM to save and query level state                                                                      | ✅     |
| Develop a user-friendly frontend dashboard                                                                      | ⏳     |
| Connect the frontend with the backend via API                                                                   | ⏳     |
| Display table of all levels                                                                                     | 📬     |
| Add a button to quickly initialize and clean a level                                                            | 📬     |
| Add dynamic configuration for each level, including objectives, tips, and solutions, displayed in the dashboard | 📬     |
| Create unit tests for frontend                                                                                  | 📬     |
| Create unit tests for backend                                                                                   | 📬     |
| Validate the product with other local Kubernetes CLIs, such as k3d, ..                                          | 📬     |

- ✅ = Ready to use
- ⏳ = In progress
- 📬 = Not yet started

## Have an idea? Notice a bug?

We'd love to hear your feedback! Feel free to log an issue on our [GitHub issues page](https://github.com/Looty/kubeground/issues).
