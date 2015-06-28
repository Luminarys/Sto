//Set the grill(s)
$(document).ready(function() {
    chinese_cartoon_pics = ["2.png", "3.png", "4.png", "5.png", "6.png", "7.png", "8.png", "9.png", "10.png", ]
    var grill = "/img/" + chinese_cartoon_pics[Math.floor(Math.random()*chinese_cartoon_pics.length)];
    $("body").css("background-image", "url(" + grill + "), url(/img/bg.png/");
});
