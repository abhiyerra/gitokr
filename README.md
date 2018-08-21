# GitSOP

Standard Operating Procedures Using Git(hub)

One of the many things about running a business is that there are Standard
Operating Procedures for many tasks. They can be things like a daily
checklist of tasks that need to be done for sales, marketing, client work,
etc. While there are tools out there that allow you to create Standard
Operating Procedures like [Process.St](https://process.st) this is a tool
built on Git(hub) and plain text so that you can easily customize it to the
way you like.

## Configuration

`.gitsop/config.json` in the root of you repo.

Config:

```
{
    "Daily Todo": {
        "cron": "* * * * *",
        "assignee": "abhiyerra",
        "files": [
            "sops/DAILY_TODO.md"
        ],
        "outputDir": "sops/output",
        "instructions": "Do this entire checklist everyday and merge it into the master branch."
    },
    "New Employee": {
        "assignee": "abhiyerra",
        "files": [
            "sops/NEW_EMPLOYEE.md"
        ],
        "outputDir": "sops/output",
        "instructions": "Do this checklist to onboard a new employee"
        "inputs": {
            "employeeName": {}
        }
    },
}
```

Notes:

 - Two different tasks
 - One is a Cron job running every minute
 - The second one must be manually invoked.

