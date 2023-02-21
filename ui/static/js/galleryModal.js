function openModal() {
    $("#galleryModal").css({
        display: "flex",
        overflow: "auto",
    });
    $("body").css("overflow", "hidden");
}

function closeModal() {
    $("#galleryModal").css("display", "none");
    $("body").css("overflow", "auto");
}

var slideIndex = 1;
showSlides(slideIndex);

function plusSlides(n) {
    showSlides((slideIndex += n));
}

function currentSlide(n) {
    showSlides((slideIndex = n));
}

function showSlides(n) {
    var i;
    var slides = $(".slides");
    if (n > slides.length) {
        slideIndex = 1;
    }
    if (n < 1) {
        slideIndex = slides.length;
    }
    slides.css("display", "none");
    $(slides[slideIndex - 1]).css("display", "flex");
}
