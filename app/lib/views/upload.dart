import 'dart:io';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class UploadPage extends StatefulWidget {
  @override
  _UploadPageState createState() => _UploadPageState();
}

class _UploadPageState extends State<UploadPage> {
  File? _selectedFile;

  Future<void> _uploadFile() async {
    if (_selectedFile == null) {
      // No file selected
      return;
    }

    var request = http.MultipartRequest(
      'POST',
      Uri.parse('/api/v1/upload'),
    );

    request.files.add(
      await http.MultipartFile.fromPath(
        'file',
        _selectedFile!.path,
      ),
    );

    var response = await request.send();

    if (response.statusCode == 200) {
      // File uploaded successfully
      print('File uploaded successfully');
    } else {
      // Error uploading file
      print('Error uploading file');
    }
  }

  Future<void> _selectFile() async {
    // TODO: Implement file selection logic
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Fraud Analysis'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton.icon(
              onPressed: _selectFile,
              icon: Icon(Icons.cloud_upload),
              label: Text('Upload'),
              style: ElevatedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(30.0),
                ),
                padding: EdgeInsets.symmetric(
                  vertical: 12.0,
                  horizontal: 24.0,
                ),
              ),
            ),
            SizedBox(height: 16.0),
            Text(
              'Upload your history here...',
              style: TextStyle(
                fontSize: 16.0,
                color: Colors.grey,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
