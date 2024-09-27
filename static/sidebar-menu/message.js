function launch_toast() {
    var toast = document.getElementById("toast")
    toast.className = "show";
    setTimeout(function(){ toast.className = toast.className.replace("show", ""); }, 5000);

    var desc = document.getElementById("desc")
    desc.className = "show";
    setTimeout(function(){ desc.className = desc.className.replace("show", ""); }, 5000);
}