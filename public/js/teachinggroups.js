// teachinggroups.js


// 
// fill the teaching group selector
// 
function fillTeachingGroupSelector(teacher_guid) {

    var $dropdown = $("#teaching-group-select");
    $($dropdown).empty();


    $.ajax({
        method: "POST",
        url: "http://localhost:1340/graphql",
        contentType: "application/json",
        data: JSON.stringify({
            query: tgQuery,
            variables: { "teacher_id": teacher_guid }
        })
    }).done(function(result) {
        // fill teaching group selector
        // console.log(result);
        $.each(result.data.my_teaching_groups, function() {
            // console.log(this);
            // get the teaching periods
            var periodDescriptor = ""
            $.each(this.TeachingGroupPeriodList.TeachingGroupPeriod, function() {
                periodDescriptor = periodDescriptor + "(" + this.DayId + ":" + this.PeriodId + ") "
            });
            $dropdown.prepend($("<option />").val(this.RefId).text(this.ShortName + " - "+ periodDescriptor));
        });
        // have to re-initialise component to render
        $($dropdown).formSelect();
    });

}


var tgQuery = `query tgQuery($teacher_id: String!) {
  my_teaching_groups(teacher_id: $teacher_id) {
    RefId
    # SchoolYear
    # LocalId
    ShortName
    LongName
    # TimeTableSubjectRefId
    # StudentList{
      # TeachingGroupStudent {
        # StudentPersonalRefId
        # StudentLocalId
        # Name {
          # Type
          # FamilyName
          # GivenName
        # }
      # }
    # }
    TeacherList{
     TeachingGroupTeacher {
       StaffPersonalRefId
       # StaffLocalId
       # Association
      Name {
        # Type
        FamilyName
        GivenName
      }
     }
    }
    TeachingGroupPeriodList{
      TeachingGroupPeriod {
        RoomNumber
        DayId
        PeriodId
      }
    }
    T
  }
}`