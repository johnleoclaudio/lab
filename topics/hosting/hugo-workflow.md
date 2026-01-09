# Hugo Workflow

## Summary:

**Local**: Write -> Build -> Push 

## On Personal Machine:

#### Setup Hugo:
```bash
# Your usual Hugo setup
hugo new site my-blog
cd my-blog
git init

# Add theme as submodule
git submodule add https://github.com/adityatelange/hugo-PaperMod.git themes/PaperMod
echo "theme = 'PaperMod'" >> hugo.toml

# Create .gitignore
echo "public/" >> .gitignore  # Don't track the build output in main branch

# Commit your source
git add .
git commit -m "Initial commit"
```

#### Build and push to separate branch:
```bash
# Build your site
hugo

# Push source to main branch
git remote add origin <your-repo-url>
git push -u origin main

# Push the public/ folder to a deploy branch
cd public
git init
git add .
git commit -m "Deploy"
git remote add origin <your-repo-url>
git push -u origin deploy --force
```

Or create a script:
```bash
#!/bin/bash

echo "Building site..."
hugo

echo "Deploying to deploy branch..."
cd public
git init
git add -A
git commit -m "Deploy $(date)"
git push -f origin main:deploy

echo "Deployment complete!"
```

## On VPS
```bash
# Clone only the deploy branch
cd /var/www/yourdomain.com
git clone -b deploy --single-branch <your-repo-url> blog

# Set permissions
sudo chown -R www-data:www-data blog/
```
## Update NGINX Config
```bash
sudo vim /etc/nginx/sites-available/yourdomain.com
```

Change the root path
```bash
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    root /var/www/yourdomain.com/blog;
    index index.html;

    location / {
        try_files $uri $uri/ =404;
    }

    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```
Test changes:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

## Deploy Script on VPS

```bash
vim ~/deploy-blog.sh
```

```bash
#!/bin/bash

BLOG_DIR="/var/www/yourdomain.com/blog"

echo "Pulling latest build..."
cd $BLOG_DIR
git pull origin deploy

echo "Setting permissions..."
sudo chown -R www-data:www-data $BLOG_DIR

echo "Deployment complete!"
```

```bash
chmod +x ~/deploy-blog.sh
```
