# Setting Up a Self-Hosted Blog on VPS with HTTPS

A step-by-step guide to deploying a website on your VPS with a custom domain and SSL certificate.

## Prerequisites

- VPS with Ubuntu (this guide uses Ubuntu)
- Domain name
- SSH access to your VPS

## Step 1: Configure DNS Settings

Point your domain to your VPS by adding DNS records in your domain registrar's control panel:

**Root domain A record:**
```
Type: A
Name: @ (or leave blank)
Value: YOUR_VPS_IP_ADDRESS
TTL: 3600
```

**WWW subdomain A record:**
```
Type: A
Name: www
Value: YOUR_VPS_IP_ADDRESS
TTL: 3600
```

> **Note:** DNS propagation can take 5 minutes to 48 hours, but typically completes within an hour.

## Step 2: Install and Configure Nginx

### Install Nginx

```bash
sudo apt update
sudo apt install nginx -y
```

### Create Site Configuration

Create a new configuration file for your domain:

```bash
sudo vim /etc/nginx/sites-available/yourdomain.com
```

Add the following configuration:

```nginx
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    
    server_name yourdomain.com www.yourdomain.com;
    
    root /var/www/yourdomain.com;
    index index.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

> **Important:** Replace `yourdomain.com` with your actual domain name.

### Disable Default Site and Enable Your Site

```bash
# Remove default nginx site
sudo rm /etc/nginx/sites-enabled/default

# Create symlink to enable your site
sudo ln -s /etc/nginx/sites-available/yourdomain.com /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

## Step 3: Create Website Directory

Set up the directory structure and create a test page:

```bash
# Create directory
sudo mkdir -p /var/www/yourdomain.com

# Set proper permissions
sudo chown -R $USER:$USER /var/www/yourdomain.com
sudo chmod -R 755 /var/www/yourdomain.com

# Create test page
echo "<h1>Hello from my VPS!</h1>" > /var/www/yourdomain.com/index.html
```

## Step 4: Configure Firewall

Ensure your firewall allows HTTP and HTTPS traffic:

```bash
# Check firewall status
sudo ufw status

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Or use the Nginx Full profile (allows both)
sudo ufw allow 'Nginx Full'

# Reload firewall
sudo ufw reload
```

> **Tip:** If you haven't enabled SSH in your firewall yet, make sure to run `sudo ufw allow ssh` before enabling ufw to avoid being locked out.

## Step 5: Install SSL Certificate

Modern browsers require HTTPS. Let's Encrypt provides free SSL certificates:

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain and install SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

Follow the interactive prompts:
1. Enter your email address (for renewal notifications)
2. Agree to the Terms of Service (Y)
3. Choose whether to share your email with EFF (optional)
4. Select option 2 to redirect HTTP traffic to HTTPS (recommended)

Certbot automatically:
- Obtains a free SSL certificate
- Configures nginx for HTTPS
- Sets up auto-renewal (certificates renew every 90 days)

### Verify Auto-Renewal

Test that automatic renewal is configured correctly:

```bash
sudo certbot renew --dry-run
```

You should see: "Congratulations, all simulated renewals succeeded"

## Step 6: Test Your Site

Visit your domain in a browser:

```
https://yourdomain.com
```

You should see your test page with a valid SSL certificate! 🎉

## Troubleshooting

### Site shows Nginx welcome page instead of your content
- Make sure you disabled the default site: `sudo rm /etc/nginx/sites-enabled/default`
- Verify your site is enabled in `/etc/nginx/sites-enabled/`
- Reload nginx: `sudo systemctl reload nginx`

### Domain not accessible
- Check DNS propagation: `nslookup yourdomain.com`
- Verify firewall rules: `sudo ufw status`
- Check nginx is running: `sudo systemctl status nginx`

### Browser shows "connection refused"
- Ensure ports 80 and 443 are open in firewall
- Check if your VPS provider has additional firewall rules in their control panel

### Check Nginx logs
```bash
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log
```

## Key Takeaways

- Configured DNS to point your domain to your VPS
- Installed and configured Nginx web server
- Set up proper firewall rules
- Secured your site with a free SSL certificate
- Configured automatic HTTPS redirect
