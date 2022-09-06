$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#button").click(function(){
        $("#message").html("Sending forgot password message!");
 
        $.ajax({
            url: "/forgot-password",
            type: "POST",
