
<p align="center">
<img src="https://fireup.live/wp-content/uploads/2018/12/fireup-logo-2.png">
</p>

# Fire-Up!

## *Putting end to project generation problem for **once** and for **all**.*

### What's the use?

Your team assigned you a new task, or you just got blessed with a new groundbreaking idea.  
The next thing you wanna do is jump right to your computer and start your project. Doesn't matter if its a software project, a word document or just a Notion markdown. 
The immediate disappointment you get is - Yo have to start with an empty canvas.

*Who likes an empty canvas? Nobody right?*

> This was heavily inspired by this [presentation](https://www.youtube.com/watch?v=iMC6QZot1YA&t=1040s), end goal is to get something similar or better.


### So its just another boilerplate generator ?

Whats the use though? There are thousand similar tools and most frameworks nowadays come with project kick-starter anyways.  
And thats exactly what the problem is, A different generator to manage everything. No customizablity and utilized only for software projects (mostly).
`fire-up` is a light tool which solves this issue, it is supposed to be used with any kind of project, technical or otherwise.

### How does it work?
`fire-up` lets you define you own template, which you can then use as a blue-print to create your project from. You can either store this blue-print on your local machine or upload to github or (s3 - coming soon ) in a private repository, from `fire-up` cli itself.  
You can also use global templates contributed by other people.

#### Some terminology:

- **Artifact**  
  A template from which root project can be generated.

- **Component**  
  A template which can be injected to a project (at a specified location)


And thats not all, `fire-up` also takes care of initialization and cleanup works like downloading dependencies etc.

### How to use?

Every artifact folder needs to have one `fire-up.json` file.
It's structure looks like this:

```javascript
{
    /* can either be 'artifact' or 'component' */
    "type":"artifact",
    /* these are strings those will be replaced globally */
    "replacements":[
        {
            "placeholder":"VAR__PACKAGE_NAME",
            "desc":"Enter package name"
        },
        {
            "placeholder":"VAR__ARTIFACT_NAME",
            "desc":"Enter project/artifacts name"
        },
        {
            "placeholder":"VAR__MAIN_APP",
            "desc":"Enter the file name for entry point."
        }
    ],
    /* there are executed in created project's directory*/
    "command_list":[
        "npm install",
        "nodemon start"
    ]
}
```

```bash
# fire-up add-local <artifact-path> --alias <artifact-alias>
fire-up add-local ../Downloads/spring-data-mongo-web-artifact --alias spring-example
```

We can then use same alias to instantiate a project out of it:

```bash
# fire-up -a <artifact-alias> -r <project-name-for-this-project>
fire-up -a spring-example -r spring-microservice
```

This will result in following output:
```bash
> fire-up -a spring-web-data-mongodb -r spring-microservice-example                                                                                 

Enter package name (VAR__PACKAGE_NAME) demo_app
Enter project/artifact name (VAR__ARTIFACT_NAME) demo_artifact
Enter the file name for entry point. (VAR__MAIN_APP) MainApp

> ls -la | grep spring-microservice-example
drwxrwxrwx  14 ray  staff   448 May 10 02:57 spring-microservice-example
```

> Remote operations:

`fire-up` searches for the artifact on 3 levels:  

1. In your local machine (inside  `~/.fire-up/`)
2. In your private artifact repo, will create one if not present with name `fire-up-artifacts`
3. In global repo.

Artifacts retrieved from remote sources are cached locally.

To add artifacts to remote repository:

```bash
fire-up add-remote ~/dev/spring-data-mongo-web-artifact --alias spring-web-data-mongodb
```
