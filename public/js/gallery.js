let slideIndex = 1;
window.onload = function() {
	showSlides(slideIndex)
}

// Next/previous controls
function plusSlides(n) {
  showSlides(slideIndex += n);
}

// Thumbnail image controls
function currentSlide(n) {
  showSlides(slideIndex = n);
}

function showSlides(n) {
  let i;
  var slides = document.getElementsByClassName("slide");
  console.log(slides.length)
  if (n > slides.length) {slideIndex = 1}
  if (n < 1) {slideIndex = slides.length}
  for (i = 0; i < slides.length; i++) {
    slides[i].style.display = "none";
  }
	slides[slideIndex - 1].style.display = "flex";
}