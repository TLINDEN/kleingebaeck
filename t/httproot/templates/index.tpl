<!DOCTYPE html>
<html lang="de" >
  <head>
    <title>Ads</title>
  </head>
  <body>
    {% for ad in ads %}
     <h2 class="text-module-begin">
        <a class="ellipsis"
           href="/s-anzeige/{{ ad.slug }}/{{ ad.id }}">{{ ad.title }}</a>
     </h2>
    {% endfor %}
  </body>
</html>

