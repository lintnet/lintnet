function(param)
  local workflow_names = std.set(std.map(
    function(data) data.value.name,
    param.combined_data
  ));

  std.filter(
    function(elem) std.length(elem.location.files) > 1,
    std.map(function(workflow_name) {
      name: 'GitHub Actions workflow name must be unique',
      location: {
        workflow_name: workflow_name,
        files: std.filterMap(
          function(data) data.value.name == workflow_name,
          function(data) data.file_path,
          param.combined_data
        ),
      },
    }, workflow_names)
  )
