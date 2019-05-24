// nias.js

// initiialise selector controls
$(document).ready(function() {
    M.AutoInit();

    $('.modal').modal();

    fillTeacherSelector();

    // 
    // respond to teacher selection
    // get the teaching groups
    // 
    $("#teacher-select").on("change", function(event) {
        $selection = $("#teacher-select").val()
        // console.log($selection);
        fillTeachingGroupSelector($selection);
        // remove any existing data
        var $progress = $("#progress-table").empty();
    });
    // 
    // respond to the teaching group selection
    // get the prgress data
    // 
    $("#teaching-group-select").on("change", function(event) {
        $selection = $("#teaching-group-select").val()
        // console.log($selection);
        fillProgressTable($selection);
    });

});

// 
// set up listeners for table-row selection
// clicking invokes extended data modal popup
// 
function initProgressTableHandler() {

	$('body').off("click"); // clear any existing handlers
    $('body').on('click', '.showmodal tr', function(event) {
        event.preventDefault();
        // alert('clicked');
        fillSyllabusView();
        showExtendedData();

    });


}


// show the extended data modal form
function showExtendedData() {
	$("#ed-modal").css("max-height", "80%");
    $('#ed-modal').modal('open');
}