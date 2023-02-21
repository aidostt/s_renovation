$(window).on("load", function () {
    $("#preloader").fadeOut("slow", function () {
        $(this).remove();
    });
});

$(window).scroll(function () {
    if ($(window).scrollTop() > 60) {
        $("#header").css({
            padding: "0.5rem 0",
            transition: "0.7s ease",
        });
    } else {
        $("#header").css({
            padding: "0.8rem 0",
            transition: "0.7s ease",
        });
    }
});

$("a").on("click", function (e) {
    if (this.hash != "") {
        // e.preventDefault();
        const hash = this.hash;

        $("html, body").animate(
            {
                scrollTop: $(hash).offset().top,
            },
            800
        );
    }
});

$(".selectPackageButton").click(function () {
    let selectedValue = $(this).data("value");
    $("#validationCustom04").val(selectedValue);
});

// gallery modal

var Offcanvas = $("#offcanvas");
var bsOffcanvas = new bootstrap.Offcanvas(Offcanvas);

$("#offcanvas a").click(function () {
    bsOffcanvas.hide();
});

// phone number mask
$(".phone-input").keyup(function () {
    var formatted = formatPhoneNumber(this.value);
    this.value = this.value.replace(/[^0-9]/g, "");
    if (this.value == "") {
        this.value = "";
    } else {
        this.value = formatted;
    }
});

function formatPhoneNumber(value) {
    var format = $(".phone-input").data("mask");
    var formatted = format;

    var i = 0;
    value.split("").forEach(function (char) {
        if (char.match(/[0-9]/)) {
            formatted = formatted.replace("#", char);
            i++;
        }
    });

    while (i < format.match(/#/g).length) {
        formatted = formatted.replace("#", "");
        i++;
    }

    return formatted;
}
