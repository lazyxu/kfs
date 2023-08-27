#getting the base image from Ubuntu 
FROM ubuntu

#Whos is maintaining the this image 
LABEL Chandra <csrkinfo@gmail.com>

#while creating docker image it will run 
RUN pat-get update 

#while creating a container this cmd will run & whiel creang docker image this will not run 
CMD ["echo", "HelloWorld...! form my first docker image"]
