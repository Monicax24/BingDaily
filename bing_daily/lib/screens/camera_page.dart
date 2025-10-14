import 'dart:io';
import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

/// Camera page for selecting and previewing images from the gallery.
class CameraPage extends StatefulWidget {
  const CameraPage({super.key});

  @override
  State<CameraPage> createState() => _CameraPageState();
}

class _CameraPageState extends State<CameraPage> {
  XFile? _selectedImage;

  /// Requests gallery permission and opens camera roll for image selection.
  Future<void> _pickImage() async {
      final ImagePicker picker = ImagePicker();
      final XFile? image = await picker.pickImage(source: ImageSource.gallery);
      // the following line is for camera access instead of gallery, but I cannot test it right now
      // final XFile? cameraImage = await picker.pickImage(source: ImageSource.camera);
      if (image != null) {
        setState(() {
          _selectedImage = image;
        });
        _showImagePreviewDialog(image);
      }
  }

  /// Displays a dialog with the selected image and Submit/Retake buttons.
  void _showImagePreviewDialog(XFile image) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: bingWhite,
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Image.file(File(image.path), height: 200, fit: BoxFit.cover),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                ElevatedButton(
                  onPressed: () {
                    Navigator.pop(context); // Close dialog
                    // Placeholder for submit action
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Image submitted')),
                    );
                    setState(() {
                      _selectedImage = null;
                    });
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: bingAccent,
                    foregroundColor: bingWhite,
                  ),
                  child: const Text(submitButtonText),
                ),
                ElevatedButton(
                  onPressed: () {
                    Navigator.pop(context); // Close dialog
                    _pickImage(); // Reopen camera roll
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: bingGreen,
                    foregroundColor: bingWhite,
                  ),
                  child: const Text(retakeButtonText),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text(appTitle), backgroundColor: bingGreen),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton(
              onPressed: _pickImage,
              style: ElevatedButton.styleFrom(
                backgroundColor: bingAccent,
                foregroundColor: bingWhite,
              ),
              child: const Text(cameraButtonText),
            ),
            if (_selectedImage != null) ...[
              const SizedBox(height: 20),
              Image.file(
                File(_selectedImage!.path),
                height: 200,
                fit: BoxFit.cover,
              ),
            ],
          ],
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
