package main

import (
	"log"
	"net/http"
)

func main() {
	shadowing()
	noShadowingSolution1()
	noShadowingSolution2()
}

func shadowing() {
	tracing := true

	var client *http.Client

	if tracing {
		// Shadowed. The scope create a new "client" variable in this scope
		client, err := createClientWithTracing()
		if err != nil {
			panic(err)
		}
		log.Println(client)

	} else {
		client, err := createDefaultClient()
		if err != nil {
			panic(err)
		}
		log.Println(client)
	}

	log.Println(client)
}

func noShadowingSolution1() {
	tracing := true

	var client *http.Client

	if tracing {
		c, err := createClientWithTracing()
		if err != nil {
			panic(err)
		}
		client = c
	} else {
		c, err := createDefaultClient()
		if err != nil {
			panic(err)
		}
		client = c
	}
	log.Println(client)
}

func noShadowingSolution2() {
	tracing := true

	var client *http.Client
	var err error

	if tracing {
		client, err = createClientWithTracing()
	} else {
		client, err = createDefaultClient()
	}

	if err != nil {
		panic(err)
	}

	log.Println(client)
}

func createClientWithTracing() (*http.Client, error) {
	return nil, nil
}

func createDefaultClient() (*http.Client, error) {
	return nil, nil
}
