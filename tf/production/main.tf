# create darkstat folder if not exists

# pull fl-data-discovery data

# run docker container with periodic patcher?
# if patch was succesful, run darkstat builder

# run nginx container to serve content?


# resource "docker_image" "darkstat_nginx" {
#   name         = "nginx:latest"
#   keep_locally = true
# }

# resource "docker_container" "darkstat_nginx" {
#   name  = "darkstat_nginx"
#   image = docker_image.darkstat_nginx.image_id

#   volumes {
#     container_path = "/usr/share/nginx/html"
#     host_path      = "/var/lib/darkstat/discovery/build"
#     read_only      = true
#   }

#   lifecycle {
#     ignore_changes = [
#       memory_swap,
#     ]
#   }
# }