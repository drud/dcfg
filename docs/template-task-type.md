# Template Task Type for App Configuration
dcfg provides a `template` Task type to help configure a web application residing within a container. `template` writes app configuration files, updates web server configuration for NGINX, and handles placement of file upload directories.

`template` currently supports the following applications. Additional apps can be supported by adding plugins, similar to tasks themselves.
- [Drupal](https://www.drupal.org)
- [WordPress](https://www.wordpress.org)

To use `template` in a task set, you must define at a minimum the template task itself and an app to configure:

```
  - name: configure drupal site
    action: template
    app: [ drupal/wordpress ]
```

## Attributes

These attributes function the same for all apps:

- `docroot`: Define the root path from which NGINX should serve a site from.
- `ignoreFiles`: Skip placement of the application's uploaded files directory.
- `databaseName`: Define the database name to connect to in the application's primary configuration file.
- `databaseUsername`: Define the database user credential in the application's primary configuration file.
- `databasePassword`: Define the database password credential in the application's primary configuration file.
- `databaseHost`: Define the database host to connect to in the application's primary configuration file.
- `databaseDriver`: Define the database driver in the application's primary configuration file.
- `databasePort`: Define the database port to connect to in the application's primary configuration file.
- `databasePrefix`: Define the database table prefix in the application's primary configuration file.

## Drupal Attributes

These attributes are available for Drupal apps:

- `siteURL`: Define the full URL of the site, including protocol. This is used to set `$base_url`.
    Allowed Value: string
    Default: empty
- `configPath`: Define a custom location for the application's primary configuration file.
    Allowed Value: string
    Default: `sites/default`
- `core`: Specify the major version release.
    Allowed Values: `7.x`, `8.x`
    Default: `7.x`
- `publicFiles`: Specify the location the public files directory should be moved to.
    Allowed Value: string
    Default: `sites/default/files`
- `privateFiles`: Specify the location the public files directory should be moved to. (Not yet implemented)
- `configSyncDir`: Used for Drupal 8, specify the location of the CMI Sync directory.
    Allowed Value: string
    Default: `/var/www/html/sync`

## WordPress Attributes

These attributes are available for WordPress apps:

- `siteURL`: Define the full URL of the site, including protocol. This is used to set `WP_SITEURL` and `WP_HOME`.
    Allowed Value: string
    Default: empty
- `configPath`: Define a custom location for the application's primary configuration file.
    Allowed Value: string
    Default: empty
- `coreDir`: If you are [Giving WordPress Its Own Directory](https://codex.wordpress.org/Giving_WordPress_Its_Own_Directory), define the directory core resides in.
    Allowed Value: string
    Default: empty
- `contentDir`: Define the content directory.
    Allowed Value: string
    Default: `wp-content`
- `uploadDir`: Define the upload directory.
    Allowed Value: string
    Default: `wp-content/uploads`
