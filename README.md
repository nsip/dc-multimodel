# dc-multimodel
simple demo server for multi-model data (SIF, xAPI and Curriculum in this case)


```
go get github.com/nsip/dc-multimodel

cd ~/go/src/github.com/nsip/dc-multimodel
go build

# run the server
./dc-multimodel

```
or just download the binary from the releases page and start with
`./dc-multimodel at your command prompt`

for best reults you should also be running an instance of the 
curriculum service, download from here.

https://github.com/nsip/dc-curriculum-service


With both services running you can log onto the web ui at

http://localhost:1340

which will show a list of teachers (SIF model), wll fetch their Teaching Groups (SIF) fetch Grading Assignements (SIF) for that group, fetch results (xAPI) for those students, link to curriculum (Curriculm json model) for task, and list Absence days (SIF) for the student.

data in graphql form is exposed at:

http://localhpst:1340/graphql


