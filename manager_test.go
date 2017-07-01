package main

import (
	"testing"

	"github.com/ashwanthkumar/gotlb/types"
	"github.com/stretchr/testify/assert"
)

const APP_ID = "/fake-app-id"

func TestManagerToCreateNewFrontendIfNotExist(t *testing.T) {
	m := NewManager()
	appInfo := createAppInfo(createAppLabels("0"))
	m.CreateNewFrontendIfNotExist(appInfo)
	f, exist := m.getFrontend(appInfo.AppId)
	assert.Equal(t, true, exist)
	f.Stop()
}

func TestManagerToAddBackendForAppShouldThrowAnErrorWhenNoFrontendIsAvailableForTheApp(t *testing.T) {
	m := NewManager()
	err := m.AddBackendForApp(createBackendInfo(APP_ID, "localhost:12345"))
	assert.Error(t, err, "Should have got an error here since frontend is not available")
}

func TestManagerToAddBackendForApps(t *testing.T) {
	m := NewManager()
	frontend := createFrontend(APP_ID, "-1", []string{"b:1", "b:2"})
	assert.Equal(t, 2, frontend.LenOfBackends())
	m.addFrontend(APP_ID, frontend)

	err := m.AddBackendForApp(createBackendInfo(APP_ID, "b:3"))
	assert.NoError(t, err)
	assert.Equal(t, 3, frontend.LenOfBackends())
}

func TestManagerToRemoveBackendForAppShouldThrowAnErrorWhenNoFrontendIsAvailableForTheApp(t *testing.T) {
	m := NewManager()
	err := m.RemoveBackendForApp(createBackendInfo(APP_ID, "localhost:12345"))
	assert.Error(t, err, "Should have got an error here since frontend is not available")
}

func TestManagerToRemoveBackendForApps(t *testing.T) {
	m := NewManager()
	frontend := createFrontend(APP_ID, "-1", []string{"b:1", "b:2"})
	assert.Equal(t, 2, frontend.LenOfBackends())
	m.addFrontend(APP_ID, frontend)

	err := m.RemoveBackendForApp(createBackendInfo(APP_ID, "b:2"))
	assert.NoError(t, err)
	assert.Equal(t, 1, frontend.LenOfBackends())
}

func createAppLabels(port string) map[string]string {
	labels := make(map[string]string)
	labels[types.TLB_PORT] = port

	return labels
}

func createAppInfo(labels map[string]string) *types.AppInfo {
	return &types.AppInfo{
		AppId:  APP_ID,
		Labels: labels,
	}
}

func createBackendInfo(appId string, backend string) *types.BackendInfo {
	return &types.BackendInfo{
		AppId: appId,
		Node:  backend,
	}
}

func createFrontend(appId, port string, backends []string) *Frontend {
	return NewFrontend(appId, port, backends)
}
