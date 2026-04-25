# Storage Provider (File Upload)

Kodia Framework provides a unified file storage system supporting multiple drivers (Local, AWS S3, Cloudflare R2) with validation, image processing, and a clean API.

---

## Configuration

### Environment Variables

```bash
# Storage Driver
STORAGE_DRIVER=local              # local, s3, cloudflare
STORAGE_MAX_FILE_SIZE=52428800    # 50MB in bytes
STORAGE_ALLOWED_MIME_TYPES=image/jpeg,image/png,application/pdf
STORAGE_ALLOWED_EXTENSIONS=.jpg,.png,.pdf
STORAGE_ENABLE_IMAGE_PROCESSING=true

# Local Storage
STORAGE_LOCAL_PATH=./storage/uploads
STORAGE_LOCAL_URL=http://localhost:8080/storage

# AWS S3
STORAGE_S3_REGION=us-east-1
STORAGE_S3_BUCKET=my-bucket
STORAGE_S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
STORAGE_S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
STORAGE_S3_CLOUDFRONT=d111111abcdef8.cloudfront.net  # Optional

# Cloudflare R2
STORAGE_CLOUDFLARE_ACCOUNT_ID=your-account-id
STORAGE_CLOUDFLARE_ACCESS_KEY=your-access-key
STORAGE_CLOUDFLARE_SECRET_KEY=your-secret-key
STORAGE_CLOUDFLARE_BUCKET=my-bucket
STORAGE_CLOUDFLARE_DOMAIN=https://my-bucket.your-domain.com
```

---

## Drivers

### Local Storage

Best for: Development, small-scale deployments, self-hosted

```bash
STORAGE_DRIVER=local
STORAGE_LOCAL_PATH=./storage/uploads
STORAGE_LOCAL_URL=http://localhost:8080/storage
```

**Features**:
- Simple setup
- Good for development
- Direct filesystem access
- No additional costs

### AWS S3

Best for: Production deployments, scalability, large files

```bash
STORAGE_DRIVER=s3
STORAGE_S3_REGION=us-east-1
STORAGE_S3_BUCKET=my-bucket
STORAGE_S3_ACCESS_KEY=AKIA...
STORAGE_S3_SECRET_KEY=...
```

**Features**:
- Highly scalable
- 99.99% durability
- Easy integration with CloudFront
- Cost-effective at scale
- Lifecycle policies for archival

### Cloudflare R2

Best for: Lower costs, global edge network, integrated CDN

```bash
STORAGE_DRIVER=cloudflare
STORAGE_CLOUDFLARE_ACCOUNT_ID=xxx
STORAGE_CLOUDFLARE_ACCESS_KEY=xxx
STORAGE_CLOUDFLARE_SECRET_KEY=xxx
STORAGE_CLOUDFLARE_BUCKET=my-bucket
```

**Features**:
- Lower costs than S3
- Global edge network
- Integrated CDN
- No egress fees
- Image optimization

---

## Basic Usage

### Single File Upload

```go
package handlers

import (
	"github.com/kodia-studio/kodia/backend/pkg/storage"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context, uploadMgr *storage.UploadManager) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}

	// Upload
	path, err := uploadMgr.Upload(c.Request.Context(), file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get URL
	url, _ := uploadMgr.GetURL(c.Request.Context(), path)

	c.JSON(200, gin.H{
		"path": path,
		"url":  url,
	})
}
```

### Image Upload with Processing

```go
func UploadAvatar(c *gin.Context, uploadMgr *storage.UploadManager) {
	file, _ := c.FormFile("avatar")

	options := &providers.ImageProcessingOptions{
		Width:           200,
		Height:          200,
		Quality:         85,
		Thumbnail:       true,
		ThumbnailWidth:  100,
		ThumbnailHeight: 100,
	}

	path, err := uploadMgr.UploadImage(c.Request.Context(), file, options)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	url, _ := uploadMgr.GetURL(c.Request.Context(), path)

	c.JSON(200, gin.H{"url": url})
}
```

### Multiple Files Upload

```go
func UploadDocuments(c *gin.Context, uploadMgr *storage.UploadManager) {
	// Get multiple files
	form, _ := c.MultipartForm()
	files := form.File["documents"]

	// Upload all
	paths, err := uploadMgr.UploadMultiple(c.Request.Context(), files)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var results []gin.H
	for _, path := range paths {
		url, _ := uploadMgr.GetURL(c.Request.Context(), path)
		results = append(results, gin.H{
			"path": path,
			"url":  url,
		})
	}

	c.JSON(200, gin.H{"files": results})
}
```

---

## File Validation

### Allowed Types

Configure allowed MIME types:

```bash
# Allow images only
STORAGE_ALLOWED_MIME_TYPES=image/jpeg,image/png,image/gif,image/webp

# Allow documents
STORAGE_ALLOWED_MIME_TYPES=application/pdf,application/msword,text/plain

# Allow multiple types
STORAGE_ALLOWED_MIME_TYPES=image/jpeg,image/png,application/pdf,text/csv
```

### Allowed Extensions

```bash
STORAGE_ALLOWED_EXTENSIONS=.jpg,.png,.gif,.pdf,.doc,.docx
```

### File Size Limits

```bash
# Set max file size (50MB)
STORAGE_MAX_FILE_SIZE=52428800

# Or per upload
options := &providers.ImageProcessingOptions{
	MaxSize: 10 * 1024 * 1024, // 10MB
}
```

---

## Image Processing

### Automatic Resizing

```go
options := &providers.ImageProcessingOptions{
	Width:   800,
	Height:  600,
	Quality: 85,
}

path, err := uploadMgr.UploadImage(ctx, file, options)
```

### Thumbnail Generation

```go
options := &providers.ImageProcessingOptions{
	Thumbnail:       true,
	ThumbnailWidth:  200,
	ThumbnailHeight: 200,
}

path, err := uploadMgr.UploadImage(ctx, file, options)
// Thumbnail stored at: path/thumbnails/{name}_thumb.ext
```

### Batch Image Processing

```go
for _, file := range files {
	options := &providers.ImageProcessingOptions{
		Width:     1920,
		Height:    1080,
		Thumbnail: true,
	}

	path, err := uploadMgr.UploadImage(ctx, file, options)
	if err != nil {
		log.Printf("Failed to upload %s: %v", file.Filename, err)
		continue
	}
}
```

---

## Advanced Usage

### Get File URL

```go
// Public URL
url, err := uploadMgr.GetURL(ctx, "2026/04/25/file.pdf")
// Result: https://bucket.s3.amazonaws.com/2026/04/25/file.pdf

// Or from config
url, err := uploadMgr.GetURL(ctx, path)
```

### File Operations

```go
// Check if file exists
exists, err := uploadMgr.Exists(ctx, path)

// Get file content
content, err := uploadMgr.GetContent(ctx, path)

// Move file
err := uploadMgr.Move(ctx, oldPath, newPath)

// Delete file
err := uploadMgr.Delete(ctx, path)

// Delete multiple
err := uploadMgr.DeleteMultiple(ctx, []string{path1, path2})
```

### Get File Info

```go
info, err := uploadMgr.GetFileInfo(ctx, path)
// Returns: FileInfo with Path, Filename, Size, MimeType, URL, Exists
```

---

## Real-World Examples

### User Avatar Upload

```go
func UpdateAvatar(c *gin.Context, user *User, uploadMgr *storage.UploadManager) {
	file, _ := c.FormFile("avatar")

	// Delete old avatar if exists
	if user.AvatarPath != "" {
		uploadMgr.Delete(c.Request.Context(), user.AvatarPath)
	}

	// Upload new avatar
	options := &providers.ImageProcessingOptions{
		Width:  200,
		Height: 200,
		Quality: 90,
		Thumbnail: true,
		ThumbnailWidth: 100,
		ThumbnailHeight: 100,
	}

	path, err := uploadMgr.UploadImage(c.Request.Context(), file, options)
	if err != nil {
		return err
	}

	// Update user
	user.AvatarPath = path
	userRepo.Update(user)

	return nil
}
```

### Product Image Gallery

```go
func UploadProductImages(c *gin.Context, product *Product, uploadMgr *storage.UploadManager) {
	form, _ := c.MultipartForm()
	files := form.File["images"]

	var imagePaths []string

	for _, file := range files {
		path, err := uploadMgr.UploadImage(c.Request.Context(), file, &providers.ImageProcessingOptions{
			Width:  1200,
			Height: 1200,
			Quality: 85,
			Thumbnail: true,
			ThumbnailWidth: 300,
			ThumbnailHeight: 300,
		})

		if err != nil {
			log.Printf("Failed to upload: %v", err)
			continue
		}

		imagePaths = append(imagePaths, path)
	}

	// Save paths to database
	product.Images = imagePaths
	productRepo.Update(product)
}
```

### Document Archive

```go
func ArchiveDocuments(c *gin.Context, archiveMgr *ArchiveManager) {
	files := form.File["documents"]

	paths, err := uploadMgr.UploadMultiple(c.Request.Context(), files)
	if err != nil {
		return err
	}

	// Store metadata
	for _, path := range paths {
		archiveMgr.SaveMetadata(path, gin.H{
			"uploaded_by": c.GetString("user_id"),
			"uploaded_at": time.Now(),
			"category": c.PostForm("category"),
		})
	}
}
```

---

## Performance Tips

✅ **Do**:
- Enable image processing for optimized storage
- Use CDN (CloudFront, R2) for global distribution
- Implement cleanup for deleted files
- Cache file URLs in database
- Use S3 lifecycle policies for archival
- Monitor storage costs

❌ **Don't**:
- Store large files directly in database
- Upload without validation
- Store sensitive files in public buckets
- Ignore file size limits
- Process images on every request

---

## Troubleshooting

### Upload Fails

1. Check file size: `STORAGE_MAX_FILE_SIZE`
2. Verify MIME type is allowed
3. Check file permissions (local storage)
4. Verify S3 credentials and bucket access
5. Check disk space (local storage)

### Slow Uploads

- Enable CDN (CloudFront, R2)
- Use multi-part upload for large files
- Check network connection
- Consider async processing

### Image Processing Issues

- Verify image format is supported
- Check memory limits
- Ensure PIL/Pillow installed
- Test with simple image first

### S3/Cloudflare Issues

- Verify credentials
- Check bucket permissions
- Ensure region is correct
- Verify bucket exists

---

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| STORAGE_DRIVER | local | Driver: local, s3, cloudflare |
| STORAGE_MAX_FILE_SIZE | 52428800 | Max file size in bytes (50MB) |
| STORAGE_ALLOWED_MIME_TYPES | All | Comma-separated allowed MIME types |
| STORAGE_ALLOWED_EXTENSIONS | All | Comma-separated allowed extensions |
| STORAGE_ENABLE_IMAGE_PROCESSING | true | Enable image processing |
| STORAGE_LOCAL_PATH | ./storage/uploads | Local storage path |
| STORAGE_LOCAL_URL | http://localhost:8080/storage | Public URL for local files |
| STORAGE_S3_REGION | us-east-1 | AWS region |
| STORAGE_S3_BUCKET | - | S3 bucket name |
| STORAGE_S3_ACCESS_KEY | - | AWS access key |
| STORAGE_S3_SECRET_KEY | - | AWS secret key |

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
