<!DOCTYPE html>
<html lang="de">
  <head>
    <title>Ad Listing</title>
  </head>
  <body>

    <div class="l-container-row">
        <div id="vap-brdcrmb" class="breadcrump">
            <a class="breadcrump-link" itemprop="url" href="/" title="Kleinanzeigen ">
                <span itemprop="title">Kleinanzeigen </span>
            </a>
            <a class="breadcrump-link" itemprop="url" href="/egal">
               <span itemprop="title">{{ category }}</span></a>
            </div>
    </div>

    {% for image in images %}
    <div class="galleryimage-element" data-ix="3">
      <img src="http://localhost:8080/api/v1/prod-ads/images/fc/{{ image.id }}?rule=$_59.JPG"/>
    </div>
    {% endfor %}

    <h1 id="viewad-title" class="boxedarticle--title" itemprop="name" data-soldlabel="Verkauft">
      {{ title }}</h1>
    <div class="boxedarticle--flex--container">
      <h2 class="boxedarticle--price" id="viewad-price">
        {{ price }}</h2>
    </div>

    <div id="viewad-extra-info" class="boxedarticle--details--full">
      <div><i class="icon icon-small icon-calendar-gray-simple"></i><span>{{ created }}</span></div>
    </div>

    <div class="splitlinebox l-container-row" id="viewad-details">
      <ul class="addetailslist">
        <li class="addetailslist--detail">
          Zustand<span class="addetailslist--detail--value" >
          {{ condition }}</span>
        </li>
      </ul>
    </div>

    <div class="l-container last-paragraph-no-margin-bottom">
      <p id="viewad-description-text" class="text-force-linebreak " itemprop="description">
        {{ text }}
      </p>
    </div>
  </body>
</html>
