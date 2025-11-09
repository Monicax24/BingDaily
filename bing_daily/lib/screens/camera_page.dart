import 'dart:io';
import 'package:bing_daily/constants.dart';
import 'package:bing_daily/providers/has_posted_provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

/// Camera page for selecting and previewing images from the gallery.
class CameraPage extends ConsumerStatefulWidget {
  const CameraPage({super.key});

  @override
  ConsumerState<CameraPage> createState() => _CameraPageState();
}

class _CameraPageState extends ConsumerState<CameraPage> {
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
                    ref.read(hasPostedProvider.notifier).state = true;
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
      body: Padding(
        padding: const EdgeInsets.only(top: 75.0),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const Padding(
                padding: EdgeInsets.symmetric(horizontal: 16.0, vertical: 38),
                child: Text(
                  'Today\'s Prompt:',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
              const Padding(
                padding: EdgeInsets.symmetric(horizontal: 28.0),
                child: Text(
                  "Share a photo of something that makes you happy!",
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: Colors.black87,
                  ),
                  textAlign: TextAlign.left,
                ),
              ),
              const SizedBox(height: 150),
              ElevatedButton(
                onPressed: _pickImage,
                style: ElevatedButton.styleFrom(
                  backgroundColor: bingGreen,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 30.0,
                    vertical: 10.0,
                  ),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8.0),
                  ),
                ),
                child: const Text(
                  cameraButtonText,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
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
      ),
      backgroundColor: bingWhite,
    );
  }
}
